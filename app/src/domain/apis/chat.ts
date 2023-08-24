import { AxiosRequestConfig } from 'axios';
import { APIERROR, normalFetch, tokenizeFetch } from './axiosInstance';

export const getChats = async () => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: 'me/chats',
  };
  try {
    const res = await normalFetch<GetChatsResp[]>(config);
    if (res.status === 200) {
      return new Chats(res.data);
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export class Chats {
  readonly chats: Array<{
    sessionName: string;
    sessionID: number;
    userID: number;
    id: number;
    content: string;
    createAt: Date;
    updateAt: Date;
    deleted: boolean;
  }> = [];
  constructor(data: Array<GetChatsResp>) {
    for (const d of data) {
      this.chats.push({
        sessionName: d.sessionName,
        sessionID: d.sessionID,
        userID: d.userID,
        id: d.id,
        content: d.content,
        createAt: new Date(d.createAt),
        updateAt: new Date(d.updateAt),
        deleted: d.deleted,
      });
    }
  }
  getIn48Chats() {
    const time = new Date();
    time.setDate(time.getDate() - 2);
    const chats: Chats['chats'] = this.chats.filter((x) => x.createAt.getTime() > time.getTime());
    chats.sort((a, b) => b.createAt.getTime() - a.createAt.getTime());
  }
}

export const getChatAt = async () => {
  const config: AxiosRequestConfig = {
    method: 'GET',
    url: 'me/chats',
  };
  try {
    const res = await normalFetch<GetChatsAtResp[]>(config);
    if (res.status === 200) {
      return new ChatAt(res.data);
    }
    return new APIERROR('', res.status);
  } catch (e) {
    console.error(e);
    return new APIERROR('', 500);
  }
};

export class ChatAt {
  readonly chats: Array<{
    id: number;
    sessionID: number;
    userID: number;
    content: string;
    createAt: Date;
    updateAt: Date;
    deleted: boolean;
  }> = [];
  constructor(data: GetChatsAtResp[]) {
    for (const d of data) {
      this.chats.push({
        id: d.id,
        sessionID: d.sessionID,
        userID: d.userID,
        content: d.content,
        createAt: new Date(d.createAt),
        updateAt: new Date(d.updateAt),
        deleted: d.deleted,
      });
    }
  }
}

export const postChatAt = async (sessionID: number, body: PostChatAtBody) => {
  const config: AxiosRequestConfig<PostChatAtBody> = {
    method: 'POST',
    url: `me/chats/${sessionID}`,
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

type GetChatsResp = {
  sessionName: string;
  sessionID: number;
  userID: number;
  id: number;
  content: string;
  createAt: string;
  updateAt: string;
  deleted: boolean;
};

type GetChatsAtResp = {
  id: number;
  sessionID: number;
  userID: number;
  content: string;
  createAt: string;
  updateAt: string;
  deleted: boolean;
};

type PostChatAtBody = {
  content: string;
};
