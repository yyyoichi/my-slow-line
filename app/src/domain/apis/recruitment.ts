import { AxiosRequestConfig } from 'axios';
import { APIERROR, normalFetch, tokenizeFetch } from './axiosInstance';

export const getRecruitments = async () => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: 'me/recruitments',
  };
  try {
    const res = await normalFetch<GetRecruitmentResp[]>(config);
    if (res.status === 200) {
      return new Recruitments(res.data);
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export class Recruitments {
  readonly recruitements: Array<{
    id: number;
    userID: number;
    message: string;
    createAt: Date;
    updateAt: Date;
    deleted: boolean;
  }> = [];
  constructor(data: GetRecruitmentResp[]) {
    for (const d of data) {
      this.recruitements.push({
        id: d.id,
        userID: d.userId,
        message: d.message,
        createAt: new Date(d.createAt),
        updateAt: new Date(d.updateAt),
        deleted: d.deleted,
      });
    }
  }
}

export const postRecruitments = async (body: PostRecruitmentBody) => {
  const config: AxiosRequestConfig<PostRecruitmentBody> = {
    method: 'POST',
    url: 'me/recruitments',
    data: body,
  };
  try {
    const res = await tokenizeFetch<null>(config);
    if (res.status === 200) {
      return true;
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export const putRecruitments = async (body: PutRecruitmentBody) => {
  const config: AxiosRequestConfig<PutRecruitmentBody> = {
    method: 'PUT',
    url: 'me/recruitments',
    data: body,
  };
  try {
    const res = await tokenizeFetch<null>(config);
    if (res.status === 200) {
      return true;
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export const getPublicRecruitmentAt = async (uuid: string) => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: `recruitments/${uuid}`,
  };
  try {
    const res = await normalFetch<GetPublicRecruitmentResp>(config);
    if (res.status === 200) {
      return res.data;
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export class PublicRecruitments {
  readonly name: string;
  readonly message: string;
  readonly uuid: string;
  constructor(data: GetPublicRecruitmentResp) {
    this.name = data.name;
    this.message = data.message;
    this.uuid = data.uuid;
  }
}

type GetRecruitmentResp = {
  id: number;
  userId: number;
  uuid: string;
  message: string;
  createAt: string;
  updateAt: string;
  deleted: boolean;
};

type PostRecruitmentBody = {
  message: string;
};

type PutRecruitmentBody = {
  uuid: string;
  message: string;
  deleted: boolean;
};

type GetPublicRecruitmentResp = {
  name: string;
  message: string;
  uuid: string;
};
