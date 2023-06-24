export const webPUshIsSupported = async () => {
  // If there is a Notification in the global space, it is considered to support the Notification API.
  if (!('Notification' in window)) {
    return false;
  }
  // If the global variable navigator has a serviceWorker property, it is considered to support service workers.
  if (!('serviceWorker' in navigator)) {
    return false;
  }
  return navigator.serviceWorker.ready.then((sw) => 'pushManager' in sw).catch(() => false);
};

export const isPwa = window.matchMedia('(display-mode: standalone)').matches;

export const hasNotification = window.Notification.permission === 'granted';
