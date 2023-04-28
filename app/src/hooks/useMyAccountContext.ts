import React from 'react';
import { MyAccount, getMe } from 'domain/apis';

export type MyAccountContextType = ReturnType<typeof useMyAccountState>;
const initContext: ReturnType<typeof useMyAccountState> = {
  myAccount: new MyAccount(),
  pullMyAccount: async () => undefined,
  isUpdated: false,
  initialized: false,
};

/**myaccount id, name, email, logined and etc.  */
export const MyAccountContext = React.createContext<MyAccountContextType>(initContext);
export const useMyAccountState = () => {
  const [initialized, setInitialized] = React.useState(false);
  const [isUpdated, setIsUpdated] = React.useState(false);
  const [state, setState] = React.useState<MyAccount>(new MyAccount());

  const pullMyAccount = async () => {
    try {
      const me = await getMe();
      setState(me);
    } catch (e) {
    } finally {
      setIsUpdated(false);
    }
  };

  // init at first
  React.useEffect(() => {
    setTimeout(async () => {
      await pullMyAccount();
      setInitialized(true);
    }, 1200);
  }, []);

  return {
    myAccount: state,
    pullMyAccount,
    isUpdated,
    initialized,
  };
};
