import React from 'react';

export type SigninFromState = {
  email: string;
  password: string;
  confirmPassword: string;
  name: string;
};

export const useSigninFormState = () => {
  const [form, setForm] = React.useState<SigninFromState>({
    email: '',
    password: '',
    confirmPassword: '',
    name: '',
  });
  const setEmail = (email: SigninFromState['email']) => setForm((pf) => ({ ...pf, email }));

  const setPassword = (password: SigninFromState['password']) => setForm((pf) => ({ ...pf, password }));

  const setConfirmPassword = (confirmPassword: SigninFromState['confirmPassword']) =>
    setForm((pf) => ({
      ...pf,
      confirmPassword,
    }));

  const setName = (name: SigninFromState['name']) => setForm((pf) => ({ ...pf, name }));

  return {
    email: form.email,
    password: form.password,
    confirmPassword: form.confirmPassword,
    name: form.name,
    setEmail,
    setPassword,
    setConfirmPassword,
    setName,
  };
};
