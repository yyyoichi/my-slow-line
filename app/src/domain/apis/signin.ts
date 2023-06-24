import { AxiosRequestConfig } from 'axios';
import { ErrorUnExpectedResponse, tokenizeFetch } from './axiosInstance';

type BasicResp = {
  jwt: string;
};

/**@return {string} jwt */
export const postSignin = async (email: string, password: string, name = '') => {
  const config: AxiosRequestConfig = {
    data: {
      email,
      password,
      name,
    },
    method: 'POST',
    url: 'signin',
  };
  try {
    const res = await tokenizeFetch<BasicResp>(config);
    if (res.status === 200) {
      return res.data.jwt;
    }
    return new Error(ErrorUnExpectedResponse);
  } catch (e) {
    throw e;
  }
};

/**@return {string} jwt */
export const postLogin = async (email: string, password: string) => {
  const config: AxiosRequestConfig = {
    data: {
      email,
      password,
    },
    method: 'POST',
    url: 'login',
  };
  try {
    const res = await tokenizeFetch<BasicResp>(config);
    if (res.status === 200) {
      return res.data.jwt;
    }
    if (res.status === 400) {
      return new Error('Could not login. Please check your email and password.');
    }
    throw new Error(ErrorUnExpectedResponse);
  } catch (e) {
    console.error(e);
    return new Error('unexpected error.');
  }
};

export const postLogout = async () => {
  const config: AxiosRequestConfig = {
    data: {},
    method: 'POST',
    url: 'me/logout',
  };
  try {
    const res = await tokenizeFetch<null>(config);
    if (res.status === 200) {
      return true;
    }
    throw new Error(ErrorUnExpectedResponse);
  } catch (e) {
    console.error(e);
    return new Error('unexpected error.');
  }
};

export const postVerificateCode = async (jwt: string, code: string) => {
  const config: AxiosRequestConfig = {
    data: {
      jwt,
      code,
    },
    method: 'POST',
    url: 'codein',
  };
  try {
    const res = await tokenizeFetch<unknown>(config);
    if (res.status === 200) {
      return true;
    }
    throw new Error(ErrorUnExpectedResponse);
  } catch (e) {
    console.error(e);
    return new Error('unexpected error.');
  }
};
