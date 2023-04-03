import { EmailValidation, ErrorEmail, ErrorName, ErrorPassword, NameValidation, PasswordValidation } from '../user';
import { NOT_NULL_ERROR } from '../validation';

describe('validation password', () => {
  const tests = [
    {
      value: '',
      error: [NOT_NULL_ERROR, ErrorPassword.Format, ErrorPassword.Length],
    },
    {
      value: '123456789',
      error: [ErrorPassword.Format],
    },
    {
      value: 'abc1',
      error: [ErrorPassword.Length],
    },
    {
      value: 'abadg',
      error: [ErrorPassword.Format, ErrorPassword.Length],
    },
  ];
  const pv = new PasswordValidation();
  const sep = ',';
  test.each(tests)("invalid-password '%s' expecte to cause '%v' errors.", ({ value, error }) => {
    const result = pv.validate(value);
    expect(result.getError(sep)).toBe(error.join(sep));
  });
});

describe('validation email', () => {
  const tests = [
    {
      value: '',
      error: [NOT_NULL_ERROR, ErrorEmail.Format],
    },
    {
      value: new Array(21).fill('a').reduce((p, x) => p + x, ''),
      error: [ErrorEmail.Format],
    },
    {
      value: new Array(51).fill('a').reduce<string>((p, x) => x + p, '@examle.com'),
      error: [ErrorEmail.Length],
    },
    {
      value: new Array(51).fill('a').reduce<string>((p, x) => x + p, ''),
      error: [ErrorEmail.Format, ErrorEmail.Length],
    },
  ];
  const ev = new EmailValidation();
  const sep = ',';
  test.each(tests)("invalid-email '%s' expecte to cause '%v' errors.", ({ value, error }) => {
    const result = ev.validate(value);
    expect(result.getError(sep)).toBe(error.join(sep));
  });
});

describe('validation name', () => {
  const tests = [
    {
      value: '',
      error: [NOT_NULL_ERROR],
    },
    {
      value: new Array(21).fill('a').reduce<string>((p, x) => p + x, ''),
      error: [ErrorName.Length],
    },
  ];
  const nv = new NameValidation();
  const sep = ',';
  test.each(tests)("invalid-name '%s' expecte to cause '%v' errors.", ({ value, error }) => {
    const result = nv.validate(value);
    expect(result.getError(sep)).toBe(error.join(sep));
  });
});
