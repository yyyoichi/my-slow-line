import { NOT_NULL_ERROR, Validation, ValidationService } from './validation';

export const ErrorVerificationCode = {
  Fromat: 'Please enter 6 numbers.',
};

export class VerificationCodeValidataion implements ValidationService<string> {
  readonly validators: ValidationService<string>['validators'] = [
    [(v) => v !== '', NOT_NULL_ERROR],
    [(v) => v.length === 6 && /\d{6}/.test(v), ErrorVerificationCode.Fromat],
  ];
  readonly InputSuggestion = 'Please enter the 6 numbers you received in your email.';
  /**@param input 6 number code */
  validate = (inputs: string) => {
    return new Validation(inputs, this.validators);
  };
}
