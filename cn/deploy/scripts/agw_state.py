from flask import Flask
from flask import Response
from prometheus_client import start_http_server, CollectorRegistry, multiprocess, generate_latest, Gauge

app = Flask(__name__)
CONTENT_TYPE_LATEST = str('text/plain; version=0.0.4; charset=utf-8')


@app.route('/metrics')
def hello_world():
    registry = CollectorRegistry()
    multiprocess.MultiProcessCollector(registry)
    return Response(generate_latest(registry), mimetype=CONTENT_TYPE_LATEST)


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000)
