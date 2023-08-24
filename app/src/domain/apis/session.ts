import { AxiosRequestConfig } from 'axios';
import { APIERROR, normalFetch, tokenizeFetch } from './axiosInstance';

export const getSessions = async () => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: '/sessions',
  };
  try {
    const res = await normalFetch<GetSessionsResp>(config);
    if (res.status === 200) {
      return new Sessions(res.data);
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export class Sessions {
  readonly id: number = 0;
  readonly name: string = '';
  readonly publicKey: string = '';
  readonly sessionStatus: SessionStatus = 'breakup';
  readonly status: ParticipantStatus = 'rejected';
  readonly createAt: Date = new Date();
  readonly updateAt: Date = new Date();
  readonly deleted: boolean = true;
  constructor(data: GetSessionsResp) {
    this.id = data.id;
    this.name = data.name;
    this.publicKey = data.publicKey;
    this.sessionStatus = data.sessionStatus;
    this.status = data.status;
    this.createAt = new Date(data.createAt);
    this.updateAt = new Date(data.updateAt);
    this.deleted = data.deleted;
  }
}

/**return sessionID */
export const postSessions = async (body: PostSessionsBody) => {
  const config: AxiosRequestConfig<PostSessionsBody> = {
    method: 'POST',
    url: 'me/sessions',
    data: body,
  };
  try {
    const res = await tokenizeFetch<PostSessionsResp>(config);
    if (res.status === 200) {
      return res.data.sessionID;
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export const getSessionAt = async (sessionID: number) => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: `me/sessions/${sessionID}`,
  };
  try {
    const res = await normalFetch<GetSessionAtResp>(config);
    if (res.status === 200) {
      return new SessionAt(res.data);
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export class SessionAt {
  readonly id: number = 0;
  readonly name: string = '';
  readonly publicKey: string = '';
  readonly sessionStatus: SessionStatus = 'breakup';
  readonly status: ParticipantStatus = 'rejected';
  readonly participants: Array<{
    id: number;
    userID: number;
    status: ParticipantStatus;
    createAt: Date;
    updateAt: Date;
    deleted: boolean;
  }> = [];
  readonly createAt: Date = new Date();
  readonly updateAt: Date = new Date();
  readonly deleted: boolean = true;
  constructor(data: GetSessionAtResp) {
    this.id = data.id;
    this.name = data.name;
    this.publicKey = data.publicKey;
    this.sessionStatus = data.sessionStatus;
    this.status = data.status;
    this.createAt = new Date(data.createAt);
    this.updateAt = new Date(data.updateAt);
    this.deleted = data.deleted;
    for (const p of data.participants) {
      this.participants.push({
        id: p.id,
        userID: p.userID,
        status: p.status,
        createAt: new Date(p.createAt),
        updateAt: new Date(p.updateAt),
        deleted: p.deleted,
      });
    }
  }
}

export const putSessionAt = async (sessionID: number, body: PutSessionAtBody) => {
  const config: AxiosRequestConfig<PutSessionAtBody> = {
    method: 'PUT',
    url: `me/sessions/${sessionID}`,
    data: body,
  };
  try {
    const res = await normalFetch<null>(config);
    if (res.status === 200) {
      return true;
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

type GetSessionsResp = {
  id: number;
  name: string;
  publicKey: string;
  sessionStatus: SessionStatus;
  status: ParticipantStatus;
  createAt: string;
  updateAt: string;
  deleted: boolean;
};

type PostSessionsBody = {
  recruitUUID: string;
  sessionName: string;
  publicKey: string;
};

type PostSessionsResp = {
  sessionID: number;
};

type GetSessionAtResp = {
  id: number;
  name: string;
  publicKey: string;
  sessionStatus: SessionStatus;
  status: ParticipantStatus;
  participants: Array<{
    id: number;
    userID: number;
    status: ParticipantStatus;
    createAt: string;
    updateAt: string;
    deleted: boolean;
  }>;
  createAt: string;
  updateAt: string;
  deleted: boolean;
};

type PutSessionAtBody = {
  sessionName: string;
};
