self.addEventListener('install', () => {
  console.log('install sw');
});
self.addEventListener('push', (event) => {
  const title = 'Test Webpush';
  const options = {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //@ts-ignore
    body: event.data.text(),
  };

  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  //@ts-ignore
  event.waitUntil(self.registration.showNotification(title, options));
});
self.addEventListener('load', (event) => {
  console.log('load', event);
});
