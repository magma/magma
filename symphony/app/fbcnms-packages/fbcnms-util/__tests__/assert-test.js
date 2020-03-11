/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import {assertType} from '../assert';

class ClassA {
  hello(): string {
    return 'hi';
  }
}
class ClassB {}
class ClassAB extends ClassA {
  hello2(): string {
    return 'hello';
  }
}

test('valid subclass', () => {
  const classA: ClassA = new ClassAB();
  classA.hello();
  const classAB = assertType(classA, ClassAB);
  classAB.hello2();
});

test('invalid subclass', () => {
  const classB: ClassB = new ClassB();
  expect(() => assertType(classB, ClassAB)).toThrowError(
    'value is not of type ClassAB',
  );
});
