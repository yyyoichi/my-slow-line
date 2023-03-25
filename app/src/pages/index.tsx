import { TokenizedFetch } from 'libs/axiosInstance';
import React from 'react';

type Person = {
  name: string;
  age: number;
};

export default function Home() {
  // const [state, setState] = React.useState<string>('');
  const onlick = () => {
    console.log('post...:');
    const tf = new TokenizedFetch();
    tf.safe({
      url: '/post',
      method: 'POST',
      data: {
        age: 23,
        name: 'hoge',
      },
    })
      .then((res) => console.log(res))
      .catch((e) => console.error(e));
  };
  React.useEffect(() => {
    onlick();
  }, []);
  return (
    <div>
      <h1>Hello world</h1>
      <p>hello</p>
      <button onClick={onlick}>fetch</button>
      {/* <div>{state}</div> */}
    </div>
  );
}
