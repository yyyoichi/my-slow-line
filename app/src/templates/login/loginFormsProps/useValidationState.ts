import React from 'react';

/**state basic validation error messages. */
export const useBasicValidationState = () => {
  const [state, setState] = React.useState<{
    email: string;
    password: string;
  }>({
    email: '',
    password: '',
  });

  /**set [email] and [password] error status message. */
  const setBasicValidationState = (email: string, password: string) =>
    setState((ps) => ({
      ...ps,
      email,
      password,
    }));

  return {
    validationState: state,
    setBasicValidationState,
  };
};

/**state validation-code validation error messages */
export const useVerificationCodeValidationState = () => {
  const [state, setState] = React.useState<string>('');
  return {
    validationState: state,
    setState,
  };
};
