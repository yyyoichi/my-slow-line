import React from 'react';
import { MyAccountContext } from 'hooks';
import { hasNotification, isPwa, webPUshIsSupported } from 'domain/web/servicesupport';
import { HomePageProps } from 'components/home';
import { getVapidPublicKey, setSubscription } from 'domain/apis/webpush';
import { isIos } from 'domain/web/device';

export const useHomePageProps = () => {
  const ac = React.useContext(MyAccountContext);
  const [notifierButtonEnable, setNotifierButtonEnable] = React.useState(true);
  const subscribeNotifier = async () => {
    try {
      setNotifierButtonEnable(false);
      const enableWebpush = await webPUshIsSupported();
      if (!enableWebpush) {
        window.alert('Your device does not support web notifications.');
        setNotifierButtonEnable(true);
        return;
      }
      const key = await getVapidPublicKey();
      if (key instanceof Error) {
        window.alert('Unexpected error occurred. Please try again.');
        setNotifierButtonEnable(true);
        return;
      }
      if (window?.Notification?.permission === 'default') {
        const result = await window?.Notification?.requestPermission();
        if (result === 'default') {
          window.alert('Reject web notifications. Please click.');
          setNotifierButtonEnable(true);
          return;
        }
      }
      if (window?.Notification?.permission === 'denied') {
        window.alert('Push notifications are blocked. Please unblock notifications from your browser settings.');
        setNotifierButtonEnable(true);
        return;
      }

      const worker = await navigator.serviceWorker.ready;
      const currentLocalSubscription = await worker.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: key,
      });
      const subscriptionJSON = currentLocalSubscription.toJSON();
      if (subscriptionJSON.endpoint == null || subscriptionJSON.keys == null) {
        window.alert(
          'The tokens issued by your browser are not yet supported, so push notifications are not available.',
        );
        setNotifierButtonEnable(true);
        return;
      }
      const resp = await setSubscription(subscriptionJSON);
      if (resp instanceof Error || !resp) {
        window.alert('Failure. Please retry.');
        setNotifierButtonEnable(true);
        return;
      }
      setNotifierButtonEnable(true);
      console.debug('success notifier');
    } catch (e) {
      console.error(e);
      window.alert('Failure. Please retry.');
      setNotifierButtonEnable(false);
    }
  };
  const reqPwa = isIos && !isPwa();
  const content = reqPwa ? 'pwa' : !ac.myAccount.has ? 'login' : !hasNotification() ? 'notification' : 'notification';
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
      tmpResult: '',
    },
  };
  return props;
};
