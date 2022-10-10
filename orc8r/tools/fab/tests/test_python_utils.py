from unittest import TestCase

from orc8r.tools.fab.python_utils import strtobool


class StrToBoolTestCase(TestCase):

    def test_true_values(self):
        self.assertTrue(strtobool("y"))
        self.assertTrue(strtobool("yes"))
        self.assertTrue(strtobool("t"))
        self.assertTrue(strtobool("true"))
        self.assertTrue(strtobool("on"))
        self.assertTrue(strtobool("1"))
        self.assertTrue(strtobool("YES"))
        self.assertTrue(strtobool("TRUE"))
        self.assertTrue(strtobool("On"))
        self.assertTrue(strtobool("Y"))
        self.assertTrue(strtobool("yEs"))

    def test_false_values(self):
        self.assertFalse(strtobool("n"))
        self.assertFalse(strtobool("no"))
        self.assertFalse(strtobool("f"))
        self.assertFalse(strtobool("false"))
        self.assertFalse(strtobool("off"))
        self.assertFalse(strtobool("0"))
        self.assertFalse(strtobool("NO"))
        self.assertFalse(strtobool("FALSE"))
        self.assertFalse(strtobool("Off"))
        self.assertFalse(strtobool("N"))
        self.assertFalse(strtobool("FaLsE"))

    def test_not_string_input(self):
        with self.assertRaises(ValueError) as context:
            strtobool(1)
        self.assertTrue("The provided value is not a string." == str(context.exception))
        with self.assertRaises(ValueError) as context:
            strtobool(0)
        self.assertTrue("The provided value is not a string." == str(context.exception))
        with self.assertRaises(ValueError) as context:
            strtobool(TestCase())
        self.assertTrue("The provided value is not a string." == str(context.exception))
        with self.assertRaises(ValueError) as context:
            strtobool(object())
        self.assertTrue("The provided value is not a string." == str(context.exception))
        with self.assertRaises(ValueError) as context:
            strtobool(True)
        self.assertTrue("The provided value is not a string." == str(context.exception))
        with self.assertRaises(ValueError) as context:
            strtobool(False)
        self.assertTrue("The provided value is not a string." == str(context.exception))

    def test_incorrect_string_input(self):
        wrong_string = "Hello"
        with self.assertRaises(ValueError) as context:
            strtobool(wrong_string)
        self.assertTrue("Value '{}' could not be converted to bool.".format(wrong_string) == str(context.exception))
        wrong_string = "treu"
        with self.assertRaises(ValueError) as context:
            strtobool(wrong_string)
        self.assertTrue("Value '{}' could not be converted to bool.".format(wrong_string) == str(context.exception))
        wrong_string = "OFFF"
        with self.assertRaises(ValueError) as context:
            strtobool(wrong_string)
        self.assertTrue("Value '{}' could not be converted to bool.".format(wrong_string) == str(context.exception))
        wrong_string = "untrue"
        with self.assertRaises(ValueError) as context:
            strtobool(wrong_string)
        self.assertTrue("Value '{}' could not be converted to bool.".format(wrong_string) == str(context.exception))
        wrong_string = "TrUtH"
        with self.assertRaises(ValueError) as context:
            strtobool(wrong_string)
        self.assertTrue("Value '{}' could not be converted to bool.".format(wrong_string) == str(context.exception))
        wrong_string = "-1"
        with self.assertRaises(ValueError) as context:
            strtobool(wrong_string)
        self.assertTrue("Value '{}' could not be converted to bool.".format(wrong_string) == str(context.exception))
        wrong_string = "truefalse"
        with self.assertRaises(ValueError) as context:
            strtobool(wrong_string)
        self.assertTrue("Value '{}' could not be converted to bool.".format(wrong_string) == str(context.exception))
