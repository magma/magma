"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from abc import ABC, abstractmethod
from collections import namedtuple
from typing import Any, Callable, Dict, List, Optional

from magma.enodebd.data_models.data_model_parameters import ParameterName

TrParam = namedtuple('TrParam', ['path', 'is_invasive', 'type', 'is_optional'])


class DataModel(ABC):
    """
    Class to represent relevant data model parameters.

    Also should contain transform functions for certain parameters that are
    represented differently in the eNodeB device than it is in Magma.

    Subclass this for each data model implementation.

    This class is effectively read-only.
    """

    def __init__(self):
        self._presence_by_param = {}

    def are_param_presences_known(self) -> bool:
        """
        True if all optional parameters' presence are known in data model
        """
        optional_params = self.get_names_of_optional_params()
        for param in optional_params:
            if param not in self._presence_by_param:
                return False
        return True

    def is_parameter_present(self, param_name: ParameterName) -> bool:
        """ Is the parameter missing from the device's data model """
        param_info = self.get_parameter(param_name)
        if param_info is None:
            return False
        if not param_info.is_optional:
            return True
        if param_name not in self._presence_by_param:
            raise KeyError(
                'Parameter presence not yet marked in data '
                'model: %s' % param_name,
            )
        return self._presence_by_param[param_name]

    def set_parameter_presence(
        self,
        param_name: ParameterName,
        is_present: bool,
    ) -> None:
        """ Mark optional parameter as either missing or not """
        self._presence_by_param[param_name] = is_present

    def get_missing_params(self) -> List[ParameterName]:
        """
        Return optional params confirmed to be missing from data model.
        NOTE: Make sure we already know which parameters are present or not
        """
        all_missing = []
        for param in self.get_names_of_optional_params():
            if self.is_parameter_present(param):
                all_missing.append(param)
        return all_missing

    def get_present_params(self) -> List[ParameterName]:
        """
        Return optional params confirmed to be present in data model.
        NOTE: Make sure we already know which parameters are present or not
        """
        all_optional = self.get_names_of_optional_params()
        all_present = self.get_parameter_names()
        for param in all_optional:
            if not self.is_parameter_present(param):
                all_present.remove(param)
        return all_present

    @classmethod
    def get_names_of_optional_params(cls) -> List[ParameterName]:
        all_optional_params = []
        for name in cls.get_parameter_names():
            if cls.get_parameter(name).is_optional:
                all_optional_params.append(name)
        return all_optional_params

    @classmethod
    def transform_for_magma(
        cls,
        param_name: ParameterName,
        enb_value: Any,
    ) -> Any:
        """
        Convert a parameter from its device specific formatting to the
        consistent format that magma understands.
        For the same parameter, different data models have their own
        idiosyncrasies. For this reason, it's important to nominalize these
        values before processing them in Magma code.

        Args:
            param_name: The parameter name
            enb_value: Native value of the parameter

        Returns:
            Returns the nominal value of the parameter that is understood
            by Magma code.
        """
        transforms = cls._get_magma_transforms()
        if param_name in transforms:
            transform_function = transforms[param_name]
            return transform_function(enb_value)
        return enb_value

    @classmethod
    def transform_for_enb(
        cls,
        param_name: ParameterName,
        magma_value: Any,
    ) -> Any:
        """
        Convert a parameter from the format that Magma understands to
        the device specific formatting.
        For the same parameter, different data models have their own
        idiosyncrasies. For this reason, it's important to nominalize these
        values before processing them in Magma code.

        Args:
            param_name: The parameter name. The transform is dependent on the
                        exact parameter.
            magma_value: Nominal value of the parameter.

        Returns:
            Returns the native value of the parameter that will be set in the
            CPE data model configuration.
        """
        transforms = cls._get_enb_transforms()
        if param_name in transforms:
            transform_function = transforms[param_name]
            return transform_function(magma_value)
        return magma_value

    @classmethod
    def get_parameter_name_from_path(
        cls,
        param_path: str,
    ) -> Optional[ParameterName]:
        """
        Args:
            param_path: Parameter path,
                eg. "Device.DeviceInfo.X_BAICELLS_COM_GPS_Status"
        Returns:
            ParameterName or None if there is no ParameterName matching
        """
        all_param_names = cls.get_parameter_names()
        numbered_param_names = cls.get_numbered_param_names()
        for _obj_name, param_name_list in numbered_param_names.items():
            all_param_names = all_param_names + param_name_list

        for param_name in all_param_names:
            param_info = cls.get_parameter(param_name)
            if param_info is not None and param_path == param_info.path:
                return param_name
        return None

    @classmethod
    @abstractmethod
    def get_parameter(cls, param_name: ParameterName) -> Optional[TrParam]:
        """
        Args:
            param_name: String of the parameter name

        Returns:
            TrParam or None if it doesn't exist
        """
        pass

    @classmethod
    @abstractmethod
    def _get_magma_transforms(
        cls,
    ) -> Dict[ParameterName, Callable[[Any], Any]]:
        """
        For the same parameter, different data models have their own
        idiosyncrasies. For this reason, it's important to nominalize these
        values before processing them in Magma code.

        Returns:
            Dictionary with key of parameter name, and value of a transform
            function taking the device-specific value of the parameter and
            returning the value in format understood by Magma.
        """
        pass

    @classmethod
    @abstractmethod
    def _get_enb_transforms(
        cls,
    ) -> Dict[ParameterName, Callable[[Any], Any]]:
        """
        For the same parameter, different data models have their own
        idiosyncrasies. For this reason, it's important to nominalize these
        values before processing them in Magma code.

        Returns:
            Dictionary with key of parameter name, and value of a transform
            function taking the nominal value of the parameter and returning
            the device-understood value.
        """
        pass

    @classmethod
    @abstractmethod
    def get_load_parameters(cls) -> List[ParameterName]:
        """
        Returns:
            List of all parameters to query when reading eNodeB state
        """
        pass

    @classmethod
    @abstractmethod
    def get_num_plmns(cls) -> int:
        """
        Returns:
            The number of PLMNs in the configuration.
        """
        pass

    @classmethod
    @abstractmethod
    def get_parameter_names(cls) -> List[ParameterName]:
        """
        Returns:
            A list of all parameter names that are neither numbered objects,
            or belonging to numbered objects
        """
        pass

    @classmethod
    @abstractmethod
    def get_numbered_param_names(
        cls,
    ) -> Dict[ParameterName, List[ParameterName]]:
        """
        Returns:
            A key for all parameters that are numbered objects, and the value
            is the list of parameters that belong to that numbered object
        """
        pass
