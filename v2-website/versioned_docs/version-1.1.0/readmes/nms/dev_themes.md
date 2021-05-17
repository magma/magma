---
id: dev_themes
title: Themes
hide_title: true
---

# Themes

Within `magmalte/app` is a `theme` directory intended to better unify the front-end design system to that found within the design files.

At the moment the `default.js` file exports 3 core design types, though more will probably follow in the future as they are needed:

## Colors
Colors are broken into subgroups (primary, secondary, state, data, and code), and then defined based off of the hex color code. Generally speaking, they should match closely with those found in the design file.

In the case you need to add a color, I used [this site](http://chir.ag/projects/name-that-color/) to generate names based off of the color hex.

## Typography
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

## Shadows
We've generated a few different elevations to leverage depending on the context in which you are designing for. In most cases, we tend to keep content within the page flat, reserving shadows for content like dialogs and modals.

