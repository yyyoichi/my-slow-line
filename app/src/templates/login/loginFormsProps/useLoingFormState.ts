import React from 'react';

export type LoginFromState = {
  email: string;
  password: string;
};

export const useLoginFormState = () => {
  const [form, setForm] = React.useState<LoginFromState>({
    email: '',
    password: '',
  });
  const setEmail = (email: LoginFromState['email']) => setForm((pf) => ({ ...pf, email }));

  const setPassword = (password: LoginFromState['password']) => setForm((pf) => ({ ...pf, password }));

  return {
    email: form.email,
    password: form.password,
    setEmail,
    setPassword,
  };
};
