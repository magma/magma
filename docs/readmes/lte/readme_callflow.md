---
id: readme_callflow
title: Building the Callflow
hide_title: true
---
# Building the Callflow

In order to visualize the attach call flow in Magma, this change adds a sequence
flow diagram. The file **Attach_call_flow_in_Magma.txt** can be uploaded to
sequencediagram.org to edit and to export the .svg. or .jpg image. The color
scheme in the diagram is as follows:

- Green: State changes
- Red: Code that crosses task boundaries or modifies *emm_context* without a function call
- Orange: Timers and notes on which function sends out the message
- Blue: Code that can be optimized, renamed or is inconsequential in this call flow
