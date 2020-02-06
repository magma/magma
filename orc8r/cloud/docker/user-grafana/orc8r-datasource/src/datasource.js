const metricNameLabel = '__name__';

// Class must implement:
//  * query(options)
//  * testDatasource()
//  * metricFindQuery(query)
//  * annotationQuery(options)
export class GenericDatasource {

  constructor(instanceSettings, $q, backendSrv, templateSrv) {
    this.type = instanceSettings.type;
    this.url = instanceSettings.url;
    this.metricsURL = this.url + "/metrics";
    this.name = instanceSettings.name;
    this.q = $q;
    this.backendSrv = backendSrv;
    this.templateSrv = templateSrv;
    this.withCredentials = instanceSettings.withCredentials;
    this.headers = {'Content-Type': 'application/json'};
    if (
      typeof instanceSettings.basicAuth === 'string' &&
      instanceSettings.basicAuth.length > 0
    ) {
      this.headers['Authorization'] = instanceSettings.basicAuth;
    }
  }

  // Query the metrics source (query_range)
  async query(options) {
    let results = [];
    for (const targetIdx in options.targets) {
      const target = options.targets[targetIdx];
      if (target.hide || target.target === "select metric") {
        continue
      }
      const result = await this.doRequest({
        url: this.buildQuery(options, target),
        data: {},
        method: 'GET',
      });
      results.push(...resultToGrafanaFormat(result))
    }
    return {data: results};
  }

  // Send GET to datasource, expect 200
  testDatasource() {
    return this.doRequest({
      url: this.url,
      method: 'GET',
    }).then(response => {
      if (response.status === 200) {
        return {
          status: 'success',
          message: 'Data source is working',
          title: 'Success',
        };
      }
    });
  }

  metricFindQuery(query) {
    // TODO: Cache results since it's unlikely for metrics series to change often
    // if empty, search all metric names, else lookahead search for metric names
    const queryMatch = query ? `{${metricNameLabel}=~"^.*${query}.*$"}` : '';

    return this.doRequest({
      url: this.metricsURL + '/series?match=' + queryMatch,
      method: 'GET',
    }).then(this.seriesToGrafanaFormat);
  }


  // for the metric query suggestions
  seriesToGrafanaFormat(result) {
    const metricNames = [...new Set(result.data.map(d => d[metricNameLabel]))];
    return metricNames.map(name => {
      return {text: name, value: name}
    });
  }

  // No concept of annotations, so leave unimplemented
  annotationQuery(options) {
    return {}
  }

  // take grafana query options and format the query URL
  buildQuery(options, target) {
    const query = this.templateSrv.replace(target.target, options.scopedVars, 'regex');

    let path = this.metricsURL + '/query_range';
    const start = options.range.from._d.valueOf() / 1000;
    const end = options.range.to._d.valueOf() / 1000;

    // Calculate step size to maximize data points without going over 11,000 (prometheus limit)
    const totalSeconds = Math.floor(end - start);
    const stepSize = Math.max(1, Math.ceil(totalSeconds/11000));

    return path
         + '?query=' + encodeURIComponent(query)
         + '&start=' + start
         + '&end=' + end
         + '&step=' + stepSize + 's';
  }

  doRequest(options) {
    options.withCredentials = this.withCredentials;
    options.headers = this.headers;
    return this.backendSrv.datasourceRequest(options);
  }
}

// convert prometheus result into a format grafana accepts
function resultToGrafanaFormat(result) {
  return result.data.data.result.map(metric => {
    return {target: buildTargetName(metric.metric),
            datapoints: formatPrometheusValues(metric.values)}
  });
}

// Add labels to the name that grafana displays for the metric
function buildTargetName(prometheusMetric) {
  const labels = Object.keys(prometheusMetric)
      .filter(key => key !== metricNameLabel)
      .map(key => `${key}=${prometheusMetric[key]}`);
  return `${prometheusMetric[metricNameLabel]}{${labels.join(',')}}`;
}

// Convert prometheus values to the format grafana expects
function formatPrometheusValues(values) {
  return values.map(val => [parseFloat(val[1]), val[0]*1000]);
}
