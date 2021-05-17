---
id: dev_components
title: Components
hide_title: true
---

# Custom Components

In order to maintain consistency and streamline the design of the NMS, we've created a few custom components to help aid in the creation of dashboards, details, and config pages. This document will serve as a point of reference for future development and should be referred back to often as the components herein continue to change and evolve over time.

## Components

### Card Title Row
The **Card Title Row** is a component you want to use anytime you are creating a title block within the app content area. **Card Title Row** takes in the following properties:

#### label `string`
Passes the label for the title row
&lt;br /&gt;

#### icon `ComponentType&lt;SvgIconExports&gt;` (Optional)
Passes in an icon to be prepended to the `label`.
&lt;br /&gt;

#### filter `() =&gt; React$Node` (Optional)
Passes in a filter on the opposite end of the row. This should be used if needing to apply something like a `date/time picker`, `edit button`, or `content filter` to the title row.
&lt;br /&gt;

##### Example

```js
 // ...

  function Filter() {
    return (
      &lt;Grid container justify="flex-end" alignItems="center" spacing={1}&gt;
        &lt;Grid item&gt;
          &lt;Text variant="body3" className={classes.dateTimeText}&gt;
            Filter By Time
          &lt;/Text&gt;
        &lt;/Grid&gt;
        &lt;Grid item&gt;
          &lt;TimeRangeSelector
            variant="outlined"
            className={classes.formControl}
            value={timeRange}
            onChange={setTimeRange}
          /&gt;
        &lt;/Grid&gt;
      &lt;/Grid&gt;
    );
  }

  return (
    &lt;&gt;
      &lt;CardTitleRow
        icon={DataUsageIcon}
        label="Gateway Check-Ins"
        filter={Filter}
      /&gt;

      ...

    &lt;/&gt;
  );
```

![image](https://user-images.githubusercontent.com/8878152/89571035-03447980-d7f5-11ea-8cc7-59e81a624846.png)

---

### Data Grid Components
The Data Grid component allows developers to quickly display information, whether it be a KPI or some other label, in whatever grid layout the design dictates. Data Grid affords developers the ability to easily add Icons, Status Indicators, Obscure Fields, and Collapse Content, all within the confines of a responsive grid in an easy-to-use format.

### Data Grid Props
`DataRows` which make up `DataGrid` take in the following properties to determine how each data entry should render and operate based on the specifics of the content you are making.
&lt;br /&gt;

#### category `string` **(Optional)**
Passes a category label for this specific data entry.
&lt;br /&gt;

#### value `string | number`
Passes a value for this specific data entry.
&lt;br /&gt;

#### unit `string` **(Optional)**
Appends a unit string to the end of a value (e.g. `'%'`)
&lt;br /&gt;

#### icon `ComponentType&lt;SvgIconExports&gt;` **(Optional)**
Passes an Icon component to be used by this specific data entry.
&lt;br /&gt;

#### obscure `boolean` **(Optional)**
Passes the obscure field to a value, making it a toggle to render the value visible.
&lt;br /&gt;

#### collapse `ComponentType | boolean` **(Optional)**
Passes an element/component as a collapsable content. Can contain additional `DataGrid` components. Can be set to false in case data may not always be available, removing the collapse from displaying.
&lt;br /&gt;

#### statusCircle `boolean` **(Optional)**
Passes a status indicator to be used by this specific data entry.
&lt;br /&gt;

#### status `boolean` **(Optional)**
Passes the state of the `statusCircle`. **True** renders a green status indicator, and **False** renders a red status indicator.
&lt;br /&gt;

#### statusInactive `boolean` **(Optional)**
Passes the state of inactive to `statusCircle` rendering a gray status indicator.
&lt;br /&gt;

#### tooltip `string` **(Optional)**
Passes a string that will be shown as a tooltip when hovering over a specific data entry. If the prop is not passed, the data entry `value` is displayed instead.
&lt;br /&gt;
&lt;br /&gt;

## Examples

### KPI's with expandable table
```js
const ran: DataRows[] = [
    [
      {
        category: 'PCI',
        value: gwInfo.cellular.ran.pci,
        statusCircle: false,
      },
      {
        category: 'eNodeB Transmit',
        value: gwInfo.cellular.ran.transmit_enabled ? 'Enabled' : 'Disabled',
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Registered eNodeBs',
        value: gwInfo.connected_enodeb_serials?.length || 0,
        collapse: &lt;EnodebsTable gwInfo={gwInfo} enbInfo={enbInfo} /&gt;,
      },
    ],
  ];

return &lt;DataGrid data={ran} /&gt;;
```
![image](https://user-images.githubusercontent.com/8878152/89434331-6bba2a80-d711-11ea-9dde-6955337e46d2.png)


### KPI's including an obscure data point
```js

  const kpiData: DataRows[] = [
    [
      {
        category: 'LTE Network Access',
        value: subscriberInfo.lte.state,
      },
    ],
    [
      {
        category: 'Data plan',
        value: dataPlan,
      },
    ],
    [
      {
        category: 'Auth Key',
        value: authKey,
        obscure: true,
      },
    ],
  ];

  return &lt;DataGrid data={kpiData} /&gt;;
```
![image](https://user-images.githubusercontent.com/8878152/89434560-ad4ad580-d711-11ea-976e-f93bd740e4f0.png)


### KPI row with parent label and icon
```js
const data: DataRows[] = [
    [
      {
        icon: CellWifiIcon,
        value: 'Gateways',
      },
      {
        category: 'Severe Events',
        value: 0,
      },
      {
        category: 'Connected',
        value: upCount || 0,
      },
      {
        category: 'Disconnected',
        value: downCount || 0,
      },
    ],
  ];

return &lt;DataGrid data={data} /&gt;;
```
![image](https://user-images.githubusercontent.com/8878152/89434865-13cff380-d712-11ea-8618-5e5bbc1db84d.png)

### KPIs with status indicators
```js
const data: DataRows[] = [
    [
      {
        category: 'Health',
        value: isGatewayHealthy(gwInfo) ? 'Good' : 'Bad',
        statusCircle: true,
        status: isGatewayHealthy(gwInfo),
      },
      {
        category: 'Last Check in',
        value: checkInTime.toLocaleString(),
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Event Aggregation',
        value: eventAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: eventAggregation,
      },
      {
        category: 'Log Aggregation',
        value: logAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: logAggregation,
      },
      {
        category: 'CPU Usage',
        value: '0',
        unit: '%',
        statusCircle: false,
      },
    ],
  ];

return &lt;DataGrid data={data} /&gt;;
```
![image](https://user-images.githubusercontent.com/8878152/89435098-627d8d80-d712-11ea-98c5-e4899fc2eb88.png)

