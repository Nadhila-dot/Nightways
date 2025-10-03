import axios, { AxiosResponse } from "axios";

const CACHE_KEY = "http_cache";
const CACHE_LIMIT = 10;
const CACHE_TTL = 60 * 1000; // 1 minute in milliseconds

interface CacheEntry {
  url: string;
  method: string;
  timestamp: number;
  data: any;
}

const getCache = (): CacheEntry[] => {
  const cached = localStorage.getItem(CACHE_KEY);
  return cached ? JSON.parse(cached) : [];
};

const setCache = (cache: CacheEntry[]) => {
  localStorage.setItem(CACHE_KEY, JSON.stringify(cache));
};

const addToCache = (url: string, method: string, data: any) => {
  const cache = getCache();
  const existingIndex = cache.findIndex(entry => entry.url === url && entry.method === method);

  if (existingIndex !== -1) {
    cache[existingIndex] = { url, method, timestamp: Date.now(), data };
  } else {
    cache.push({ url, method, timestamp: Date.now(), data });
    if (cache.length > CACHE_LIMIT) {
      cache.shift(); // Remove oldest
    }
  }

  setCache(cache);
  console.log(`🔵 CACHE STORED: ${method.toUpperCase()} ${url}`);
};

const getFromCache = (url: string, method: string): any | null => {
  const cache = getCache();
  const entry = cache.find(entry => entry.url === url && entry.method === method);

  if (entry && Date.now() - entry.timestamp < CACHE_TTL) {
    console.log(`🟢 CACHE HIT: ${method.toUpperCase()} ${url}`);
    return entry.data;
  }

  console.log(`🟠 CACHE MISS: ${method.toUpperCase()} ${url}`);
  return null;
};

const http = axios.create({
  baseURL: "/",
  headers: {
    "Content-Type": "application/json",
  },
});

http.interceptors.request.use(async (config) => {
  const session = localStorage.getItem("session");
  if (session) {
    config.headers["Authorization"] = `Bearer ${session}`;
  }
  
  if (config.method?.toLowerCase() === 'get') {
    const cachedData = getFromCache(config.url || "", 'get');
    if (cachedData) {
      const error: any = new Error('CACHED_RESPONSE');
      error.config = config;
      error.cachedData = cachedData;
      throw error;
    }
  }
  
  console.log(`🔴 HTTP REQUEST: ${config.method?.toUpperCase()} ${config.url}`);
  return config;
});

http.interceptors.response.use(
  (response) => {
    if (response.config.method?.toLowerCase() === 'get') {
      addToCache(response.config.url || "", 'get', response.data);
    }
    return response;
  }, 
  async (error) => {
    if (error.message === 'CACHED_RESPONSE' && error.cachedData) {
      console.log(`✅ SERVING FROM CACHE: ${error.config.method?.toUpperCase()} ${error.config.url}`);
      return {
        data: error.cachedData,
        status: 200,
        statusText: "OK",
        headers: {
          "x-nadhi-cache": "Client",
        },
        config: error.config,
        cached: true
      };
    }
    
    if (error.config && error.config.method?.toLowerCase() === 'get') {
      const cachedData = getFromCache(error.config.url || "", 'get');
      if (cachedData) {
        console.log(`⚠️ ERROR FALLBACK TO CACHE: ${error.config.method?.toUpperCase()} ${error.config.url}`);
        return {
          data: cachedData,
          status: 200,
          statusText: "OK",
          headers: {
            "x-nadhi-cache": "Client-Fallback",
          },
          config: error.config,
          cached: true
        };
      }
    }
    
    return Promise.reject(error);
  }
);

export default http;