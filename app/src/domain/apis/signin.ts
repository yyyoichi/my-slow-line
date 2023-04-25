import { AxiosRequestConfig } from 'axios';
import { ErrorUnExpectedResponse, tokenizeFetch } from './axiosInstance';

/**@return {number} userId */
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
    const res = await tokenizeFetch<number>(config);
    if (res.status === 200) {
      return res.data;
    }
    throw new Error(ErrorUnExpectedResponse);
  } catch (e) {
    throw e;
  }
};

/**@return {number} userId */
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
    const res = await tokenizeFetch<number>(config);
    if (res.status === 200) {
      return res.data;
    }
    if (res.status === 400) {
      throw new Error('Could not login. Please check your email and password.');
    }
    throw new Error(ErrorUnExpectedResponse);
  } catch (e) {
    throw e;
  }
};

/**@param id userID */
export const postVerificateCode = async (id: number, code: string) => {
  const config: AxiosRequestConfig = {
    data: {
      id,
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
    throw e;
  }
};
