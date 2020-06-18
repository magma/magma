/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import SvgIcon from '@material-ui/core/SvgIcon';

type Props = {
  className?: string,
};

const ArrowRightCircle = ({className}: Props) => (
  <SvgIcon>
    <defs>
      <rect x="29" y="415" width="20" height="20" id="rect-1" />
    </defs>
    <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
      <g id="Select-state" transform="translate(-29.000000, -415.000000)">
        <g id="glyph-arrow-right-circle">
          <image
            x="29"
            y="415"
            width="20"
            height="20"
            className={className}
            href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACgAAAAoCAYAAACM/rhtAAAABGdBTUEAALGPC/xhBQAAAYRpQ0NQSUNDIFByb2ZpbGUAACiRlZE/SMNAFMa/VKSD9c9QRJwCSqcqtYIoiNB2EKFDLR1UdEgv6R+oaUjPqugo6FhwEF0Uuzg461pw1UlRBHFxcS+6aInvkmoUUfGD8H559+7eu+8Az5FiGAVPF7CoczM5GZVnZudk7wO8aIMthZWMSCIRF/wev+r5GpKIlwPirO/rv8qnaiUGSH5izgyTE28R9y1zQ/Axsd+koYhrgrMOXwlOO/xo16SSMeJX4o4yy9Jej484pKt5nXiceILlFJV4nTiY/lSf/cTOPLZ6Y0ohnzYVrqmysCZWLBTNkqEw7Z+X/EtcW+EixorGqpnP5rgcISc1eUpng0E5HBoaA8S7ONX1adtvqfvczRUPgNEnoKXi5tI7wOkm0HPr5vr3gc4N4OSMLZnlZnsJ1N3xrPkvfRB+YsdXWyNA9QZIrQHxC2B3DwhkqM88kGinPE3PF9zPUSkzHHbO8kWB1nvLqgcA7zbQqFjWy6FlNap0nzugpr8BF6NvkBl0THQAAANbSURBVFgJzZk/a9ZQGMVTUVCsCg4d6mh1cussb9xEC1VXhw5+A79A4XXt1+giKDgoLlL8APUP3RQhg9YO3argINbzi3leniap3pvcSA6c5ia5z3lO7s39k7dzWXdcVuiquCwuOqqY7Tpuq/xM/CQOjgVlmIo74mEkiSEWjeSYl+K6eCDGGqvXRwMtNJPgnlT2xHqib7r2RHwg5uJVkaSQci5yjzrUrcejiXZnzCmSJ/0levE3Or8rnhZDQV1iiPVaaJODXFE4o9qPRS9W6Py+GC2mGAOxaBSi1yYXOYOASN3cC127EBQdVgktNOsmgx6eJveBGzo/IaYGmmj7XOT+K3hp/TuHwNDwJsl97MCZ100/WumCIVqu/sDk8N2NB7w04Lu20N2U71wjWe0CuQrRurvR1Qu66SdhRlpXnFXgK/GdeD1ChJxmEC94muGRSnaTuSpoNM2ijxZWnBYT9OTo7WPPyOnnyamvuaMTM3jH3+hQ5v357PRiTDKZmw88lWBXYhcRi1kh/ig0/y7pUheT5PbLIt6yh6IZZN1Mha4m8WB+8JZtugss7ilxRWKxLYkHM7jJHMRm05B6U/lRwjfEL1UCRvhzcVKdtx28h8WTquENshNuwzVdXBVPtd0MuPZUddbE8yImX4o3xddiHd5D6c3Pf20zOA+xL1qzpzrS9W0rFR4sxwEVOBkrDmmdr+K5yiFN+qEq2+GnCrnYp4svKn5NpIvBD5GVg81BHf6Vw1u2JVqT5lxIjNiRnCu/+dmii/1LWU6MCQ1ijga4VGl+1/G22DY4qiqZ97CLwW27o+MtV+5bXJJArDlyeg+lNxxbk45yqcO13yywYPcB0wRTiH/oSaBg62aB2KkTZMvD1qcrVhTYxRw537pYPM0wxIb1vdSTbVhxui7akxcqsw3/XyBXIVp+vDTAu8MHi1XiQ4ZRPjTIQS7Liwe8tIJPvtF+dppjmteehuOGOERLoom2z0Xuf4LRNOqfPniCUf94ZE1MS9Lk/p2kO5gnmVD5wAkFdYkh1ncp2uToM++Wv5X40W0JWBb5wOEbIhf50ZLRBynnIveo47/ULB5NBmUSkJQn9btvSxR7RAMtNJODFWcq+rU71CAxxKIRjD59zy6IXfayyC7YqGK6f0P8BnIvjbCPSBPIAAAAAElFTkSuQmCC"
          />
          <use fill="#B3B3B3" fill-rule="evenodd" href="#rect-1" />
        </g>
      </g>
    </g>
  </SvgIcon>
);

export default ArrowRightCircle;
