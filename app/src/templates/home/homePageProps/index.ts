import React from 'react';
import { MyAccountContext } from 'hooks';
import { hasNotification, isPwa, webPUshIsSupported } from 'domain/web/servicesupport';
import { HomePageProps } from 'components/home';
import { getVapidPublicKey } from 'domain/apis/webpush';
import { isIos } from 'domain/web/device';

export const useHomePageProps = () => {
  const ac = React.useContext(MyAccountContext);
  const [notifierButtonEnable, setNotifierButtonEnable] = React.useState(true);
  const [subscription, setSubscription] = React.useState('');
  React.useEffect(() => {
    webPUshIsSupported().then((res) => {
      console.log('push support', res);
    });
  }, []);
  const subscribeNotifier = async () => {
    setNotifierButtonEnable(false);
    const enableWebpush = await webPUshIsSupported();
    if (!enableWebpush) {
      window.alert('Your device is not support web notifications.');
      return setNotifierButtonEnable(true);
    }
    const key = await getVapidPublicKey();
    if (key instanceof Error) {
      window.alert('Unexpected error occured. Please try agein.');
      return setNotifierButtonEnable(true);
    }
    console.log(key);
    if (window.Notification.permission === 'default') {
      const result = await window.Notification.requestPermission();
      if (result === 'default') {
        window.alert('Reject web notifications. Please click');
        return setNotifierButtonEnable(true);
      }
    }
    if (window.Notification.permission === 'denied') {
      window.alert('Push notifications are blocked. Please unblock notifications from your browser settings.');
      return setNotifierButtonEnable(true);
    }

    const currentLocalSubscription = await navigator.serviceWorker.ready.then((worker) =>
      worker.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: key,
      }),
    );
    const subscriptionJSON = currentLocalSubscription.toJSON();
    console.log(subscriptionJSON);
    if (subscriptionJSON.endpoint == null || subscriptionJSON.keys == null) {
      window.alert('The tokens issued by your browser are not yet supported, so push notifications are not available.');
      return setNotifierButtonEnable(true);
    }
    setSubscription(JSON.stringify(subscriptionJSON));
  };
  const reqPwa = isIos && !isPwa;
  const content = !ac.myAccount.has ? 'login' : reqPwa ? 'pwa' : !hasNotification ? 'notification' : 'notification';
  console.log(content);
  const props: HomePageProps = {
    switchContext: {
      content,
    },
    notification: {
      notificationButton: {
        active: notifierButtonEnable,
        onClick: () => {
          subscribeNotifier();
        },
      },
      tmpResult: subscription,
    },
  };
  return props;
};
