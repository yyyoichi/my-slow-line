type Input = string | number;

export const NOT_NULL_ERROR = 'Please enter.';

/**interface for components. at first initialize, and perform validation when submit event*/
export interface ValidationService<T extends Input | Array<Input>> {
  readonly InputSuggestion: string;

  /**a element for performing validation and prepare the error message */
  readonly validators: [
    /**validate function expected to return boolean */
    (v: T) => boolean,
    /**error message when invalid */
    string,
  ][];

  validate: (input: T) => Validation<T>;
}

/**holds validation result end errors message. */
export class Validation<T extends Input | Array<Input>> {
  readonly isValid: boolean;
  readonly errors: string[];
  /**has the result of validation check and sujest error message for user.
   * @param target the value you wont to validate
   * @param validators  validations
   */
  constructor(target: T, validators: ValidationService<T>['validators']) {
    const errs: string[] = [];

    // check all validation
    for (const [validate, err] of validators) {
      if (validate(target)) continue;
      errs.push(err);
    }
    this.errors = errs;
    this.isValid = errs.length === 0;
  }
  /**get error message to string.
   * @param separator A string used to separate one element of the array.
   */
  getError(separator = ', ') {
    return this.errors.join(separator);
  }
}
