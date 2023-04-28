import React from 'react';

/**state basic validation error messages. */
export const useBasicValidationState = () => {
  const [state, setState] = React.useState<{
    email: string;
    password: string;
    confirmPassword: string;
    name: string;
  }>({
    email: '',
    password: '',
    confirmPassword: '',
    name: '',
  });

  /**set [email], [password] and [confirmPassword] error status message. */
  const setBasicValidationState = (email: string, password: string, confirmPassword: string) =>
    setState((ps) => ({
      ...ps,
      email,
      password,
      confirmPassword,
    }));
  const setNameState = (name: string) => setState((ps) => ({ ...ps, name }));

  return {
    validationState: state,
    setBasicValidationState,
    setNameState,
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
