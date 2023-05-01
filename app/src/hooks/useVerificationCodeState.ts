import React from 'react';

type CodesState = {
  jwt: string;
  code: string;
};

export const useVerificationCodeState = () => {
  const [state, setState] = React.useState<CodesState>({
    jwt: '',
    code: '',
  });
  const setJwt = (jwt: CodesState['jwt']) => setState((ps) => ({ ...ps, jwt }));
  const setCode = (code: CodesState['code']) => setState((ps) => ({ ...ps, code }));

  return {
    code: state.code,
    jwt: state.jwt,
    setJwt,
    setCode,
  };
};
