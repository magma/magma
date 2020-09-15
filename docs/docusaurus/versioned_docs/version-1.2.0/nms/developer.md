---
id: version-1.2.0-nms_developer
title: Quick-Start Guide to NMS Development
hide_title: true
original_id: nms_developer
---
# Spacing Guidelines
Within the context of the NMS, it's important to have consistency within the design so that all pages feel unified. As such, most spacing throughout the app follows an 8px scaling factor as shown below:

![image](https://user-images.githubusercontent.com/8878152/89579208-5e7c6900-d801-11ea-8456-bef4313cde40.png)
![image](https://user-images.githubusercontent.com/8878152/89579597-07c35f00-d802-11ea-9396-f056c0e82b4a.png)
![image](https://user-images.githubusercontent.com/8878152/89579632-1c075c00-d802-11ea-88c4-0dc3179cd2c0.png)

To help facilitate this better, we leverage [Material-UI's](https://material-ui.com/customization/spacing/#spacing) `theme.spacing()` helper which too uses an 8px scaling factor.

```js
const theme = createMuiTheme();

theme.spacing(0.5) // = 8 * 0.5 (4px)
theme.spacing(1) // = 8 * 1 (8px)
theme.spacing(2) // = 8 * 2 (16px)
theme.spacing(3) // = 8 * 3 (24px)
theme.spacing(4) // = 8 * 4 (32px)
theme.spacing(5) // = 8 * 5 (40px)
```

With this in mind, always try and leverage the scaling system when building out components rather than using static `px` values. Reason being, in the case the scaling factor is ever changed in the future, it will automatically update across all sizing.

# Local Theming
Within `magmalte/app` is a `theme` directory intended to better unify the front-end design system to that found within the design files.

At the moment the `default.js` file exports 3 core design types, though more will probably follow in the future as they are needed:

#### Colors
Colors are broken into subgroups (primary, secondary, state, data, and code), and then defined based off of the hex color code. Generally speaking, they should match closely with those found in the design file.

In the case you need to add a color, I used [this site](http://chir.ag/projects/name-that-color/) to generate names based off of the color hex.

#### Typography
We've adapted the Symphony styles to better reflect the typography found in the design file, with the following possible variants:

* h1
* h2
* h3
* h4 
* h5
* subtitle1
* subtitle2
* body1
* body2
* body3
* code
* button
* caption
* overline

#### Shadows
We've generated a few different elevations to leverage depending on the context in which you are designing for. In most cases, we tend to keep content within the page flat, reserving shadows for content like dialogs and modals.

# Custom Development Components

In order to maintain consistency and streamline the design of the NMS, we've created a few custom components to help aid in the creation of dashboards, details, and config pages. This document will serve as a point of reference for future development and should be referred back to often as the components herein continue to change and evolve over time.

## Components

### Card Title Row
The **Card Title Row** is a component you want to use anytime you are creating a title block within the app content area. **Card Title Row** takes in the following properties:

#### label `string`
Passes the label for the title row
<br />

#### icon `ComponentType<SvgIconExports>` (Optional)
Passes in an icon to be prepended to the `label`.
<br />

#### filter `() => React$Node` (Optional)
Passes in a filter on the opposite end of the row. This should be used if needing to apply something like a `date/time picker`, `edit button`, or `content filter` to the title row.
<br />

##### Example

```js

 ...

  function Filter() {
    return (
      <Grid container justify="flex-end" alignItems="center" spacing={1}>
        <Grid item>
          <Text variant="body3" className={classes.dateTimeText}>
            Filter By Time
          </Text>
        </Grid>
        <Grid item>
          <TimeRangeSelector
            variant="outlined"
            className={classes.formControl}
            value={timeRange}
            onChange={setTimeRange}
          />
        </Grid>
      </Grid>
    );
  }

  return (
    <>
      <CardTitleRow
        icon={DataUsageIcon}
        label="Gateway Check-Ins"
        filter={Filter}
      />

      ...

    </>
  );
```

![image](https://user-images.githubusercontent.com/8878152/89571035-03447980-d7f5-11ea-8cc7-59e81a624846.png)

---

### Data Grid Components
The Data Grid component allows developers to quickly display information, whether it be a KPI or some other label, in whatever grid layout the design dictates. Data Grid affords developers the ability to easily add Icons, Status Indicators, Obscure Fields, and Collapse Content, all within the confines of a responsive grid in an easy-to-use format.

### Data Grid Props
`DataRows` which make up `DataGrid` take in the following properties to determine how each data entry should render and operate based on the specifics of the content you are making.
<br />

#### category `string` **(Optional)** 
Passes a category label for this specific data entry.
<br />

#### value `string | number`
Passes a value for this specific data entry.
<br />

#### unit `string` **(Optional)**
Appends a unit string to the end of a value (e.g. `'%'`)
<br />

#### icon `ComponentType<SvgIconExports>` **(Optional)**
Passes an Icon component to be used by this specific data entry.
<br />

#### obscure `boolean` **(Optional)**
Passes the obscure field to a value, making it a toggle to render the value visible.
<br />

#### collapse `ComponentType | boolean` **(Optional)**
Passes an element/component as a collapsable content. Can contain additional `DataGrid` components. Can be set to false in case data may not always be available, removing the collapse from displaying.
<br />

#### statusCircle `boolean` **(Optional)**
Passes a status indicator to be used by this specific data entry.
<br />

#### status `boolean` **(Optional)**
Passes the state of the `statusCircle`. **True** renders a green status indicator, and **False** renders a red status indicator.
<br />

#### statusInactive `boolean` **(Optional)**
Passes the state of inactive to `statusCircle` rendering a gray status indicator. 
<br />

#### tooltip `string` **(Optional)**
Passes a string that will be shown as a tooltip when hovering over a specific data entry. If the prop is not passed, the data entry `value` is displayed instead.
<br />
<br />

### Examples

#### KPI's with expandable table
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
        collapse: <EnodebsTable gwInfo={gwInfo} enbInfo={enbInfo} />,
      },
    ],
  ];

return <DataGrid data={ran} />;
```
![image](https://user-images.githubusercontent.com/8878152/89434331-6bba2a80-d711-11ea-9dde-6955337e46d2.png)


#### KPI's including an obscure data point
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

  return <DataGrid data={kpiData} />;
```
![image](https://user-images.githubusercontent.com/8878152/89434560-ad4ad580-d711-11ea-976e-f93bd740e4f0.png)


#### KPI row with parent label and icon
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

return <DataGrid data={data} />;
```
![image](https://user-images.githubusercontent.com/8878152/89434865-13cff380-d712-11ea-8618-5e5bbc1db84d.png)

#### KPIs with status indicators
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

return <DataGrid data={data} />;
```
![image](https://user-images.githubusercontent.com/8878152/89435098-627d8d80-d712-11ea-98c5-e4899fc2eb88.png)

