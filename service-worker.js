console.log('service worker getting out of bed and having a coffee');

self.addEventListener("install", function(event) {
   self.skipWaiting();

   event.waitUntil(
      caches.open(CacheName)
      .then(cache => {
         return cache.addAll(CacheUrls);
      })
   );
});

self.addEventListener('activate', (event) => {
    // newest version of cache and only one that is allowed
    const cacheAllowList = ['VentureOut-cache-v13'];
  
    // Get all the caches that a currently active.
    event.waitUntil(caches.keys().then((keys) => {
      // Delete all caches that are not the allowed cache
      return Promise.all(keys.map((key) => {
        if (!cacheAllowList.includes(key)) {
          return caches.delete(key);
        }
      }));
    }));
  });

self.addEventListener('fetch', function(event) {
    // console.log('service worker fetch event', event);
    event.respondWith(
        caches.match(event.request).then( function(cacheResponse) {
            if (cacheResponse) {
                return cacheResponse;
            }
            return fetch(event.request);
        })
    );
});

const CacheName = 'VentureOut-cache-v13';
const CacheUrls = [
    './',
    '/static/css/main.css', 
    '/static/css/mdb.dark.min.css',
    '/static/css/mdb.dark.rtl.min.css',
    '/static/css/mdb.min.css',
    '/static/css/mdb.rtl.min.css',
    '/static/img/favicon.ico',
    '/static/img/apple-touch-icon.png',
    '/static/img/favicon-16x16.png',
    '/static/img/favicon-32x32.png',
    '/static/img/logo.png',
    '/static/img/navlogo.png',
    '/static/img/maskable_icon_x72.png',
    '/static/img/maskable_icon_x96.png',
    '/static/img/maskable_icon_x128.png',
    '/static/img/maskable_icon_x192.png',
    '/static/img/maskable_icon_x384.png',
    '/static/img/maskable_icon_x512.png',
    '/static/img/hwssc.png',
    '/static/img/heritage.jpg',
    '/static/img/discovery.png',
    '/static/img/fitness.png',
    '/static/img/sasc.jpg',
    '/static/js/html5-qrcode.min.js',
    '/static/js/mdb.min.js',
    'manifest.webmanifest',
    'service-worker.js',
    '/static/templates/addevent.html',
    '/static/templates/ar.html',
    '/static/templates/footer.html',
    '/static/templates/header.html',
    '/static/templates/index.html',
    '/static/templates/calender.html',
    '/static/templates/infomodal.html',
    '/static/templates/map.html',
    '/static/templates/navbar.html',
    '/static/templates/navbarmobile.html',
    '/static/templates/scan.html',
    '/static/templates/scripts.html',
    '/static/templates/settings.html'
];