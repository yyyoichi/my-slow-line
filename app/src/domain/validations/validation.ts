type Input = string | number;

export const NOT_NULL_ERROR = 'Please enter.';

/**interface for components. at first initialize, and perform validation when submit event*/
export interface ValidationInput<T extends Input> {
  readonly InputSuggestion: string;
  validate: (input: T) => Validation<T>;
}

/**a element for performing validation and prepare the error message */
export type Validator<T extends Input> = [
  /**validate function expected to return boolean */
  (v: T) => boolean,
  /**error message when invalid */
  string,
];

/**holds validation result end errors message. */
export class Validation<T extends Input> {
  readonly isValid: boolean;
  readonly errors: string[];
  /**has the result of validation check and sujest error message for user.
   * @param target the value you wont to validate
   * @param validators  validations
   */
  constructor(target: T, validators: Validator<T>[]) {
    const errs: string[] = [];

    // check all validation
    for (const v of validators) {
      if (v[0](target)) continue;
      errs.push(v[1]);
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
