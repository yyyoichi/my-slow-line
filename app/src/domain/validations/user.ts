import { NOT_NULL_ERROR, Validation, ValidationInput, Validator } from './validation';

export const ErrorPassword = {
  Format: 'Please use single-byte alphanumeric characters.',
  Length: 'Please enter between 8 and 24 characters',
};

export const ErrorEmail = {
  Format: 'Please enter a valid Email address.',
  Length: 'Please enter within 50 characters.',
};

export const ErrorName = {
  Length: 'Please enter within 20 characters.',
};

export class PasswordValidation implements ValidationInput<string> {
  private static validators: Validator<string>[] = [
    [(v) => v !== '', NOT_NULL_ERROR],
    [(v) => /[0-9].*[a-zA-Z]|[a-zA-Z].*[0-9]/.test(v), ErrorPassword.Format],
    [(v) => 8 <= v.length && v.length <= 24, ErrorPassword.Length],
  ];
  readonly InputSuggestion = 'Half-width alphanumeric characters, 8 to 24 characters';
  validate = (input: string) => {
    return new Validation(input, PasswordValidation.validators);
  };
}

export class EmailValidation implements ValidationInput<string> {
  private static validators: Validator<string>[] = [
    [(v) => v !== '', NOT_NULL_ERROR],
    [(v) => /.+@.+/.test(v), ErrorEmail.Format],
    [(v) => v.length <= 50, ErrorEmail.Length],
  ];
  readonly InputSuggestion = 'Email address within 50 single-byte characters.';
  /**expected email is not null, max 50 and inclued '*@*' */
  validate = (input: string) => {
    return new Validation(input, EmailValidation.validators);
  };
}

export class NameValidation implements ValidationInput<string> {
  private static validators: Validator<string>[] = [
    [(v) => v !== '', NOT_NULL_ERROR],
    [(v) => v.length <= 20, ErrorName.Length],
  ];
  readonly InputSuggestion = 'User Name within 20 characters.';
  /**expected email is not null, max 50 and inclued '*@*' */
  validate = (input: string) => {
    return new Validation(input, NameValidation.validators);
  };
}
