"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from abc import abstractmethod, ABC
from magma.enodebd.device_config.enodeb_configuration import \
    EnodebConfiguration


class EnodebConfigurationPostProcessor(ABC):
    """
    Overrides the desired configuration for the eNodeB, with subclass per
    device/sw-version that requires non-standard configuration behavior.
    """

    @abstractmethod
    def postprocess(self, desired_cfg: EnodebConfiguration) -> None:
        """
        Implementation of function which overrides the desired configuration
        for the eNodeB
        """
        pass
