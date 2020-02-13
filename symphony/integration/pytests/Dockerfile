# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM python:3.6

COPY cli/ ./
RUN python setup.py -q install

# Switch to intergration pytests directory
WORKDIR /src/symphony/integration/pytests

CMD ["python", "runner.py"]
