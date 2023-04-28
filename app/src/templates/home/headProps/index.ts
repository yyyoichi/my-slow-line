import React from 'react';
import { MyAccountContext } from 'hooks';
import { HeadProps } from 'components/home';

export const useHeadProps = () => {
  const ac = React.useContext(MyAccountContext);
  const props: HeadProps = {
    isLogined: ac.myAccount.has,
    account: {
      name: ac.myAccount.name,
    },
  };
  return props;
};
