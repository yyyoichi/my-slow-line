import React from 'react';
import { MyAccountContext } from 'hooks';
import { HeadProps } from 'components/head';
import { postLogout } from 'domain/apis/signin';

export const useHeadProps = () => {
  const ac = React.useContext(MyAccountContext);
  const props: HeadProps = {
    isLogined: ac.myAccount.has,
    user: {
      icon: {
        children: ac.myAccount.name[0],
      },
      logout: {
        onClick: () => {
          postLogout();
          ac.clearMyAccount();
        },
      },
    },
  };
  return props;
};
