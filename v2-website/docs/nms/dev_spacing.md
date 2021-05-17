---
id: dev_spacing
title: Spacing Guidelines
hide_title: true
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

