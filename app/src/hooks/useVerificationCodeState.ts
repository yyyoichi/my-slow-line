import React from 'react';

type CodesState = {
  userId: number;
  code: string;
};

export const useVerificationCodeState = () => {
  const [state, setState] = React.useState<CodesState>({
    userId: 0,
    code: '',
  });
  const setUserId = (userId: CodesState['userId']) => setState((ps) => ({ ...ps, userId }));
  const setCode = (code: CodesState['code']) => setState((ps) => ({ ...ps, code }));

  return {
    code: state.code,
    userId: state.userId,
    setUserId,
    setCode,
  };
};
