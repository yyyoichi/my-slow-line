import React from 'react';

export default function Home() {
  const [state, setState] = React.useState<string>('');
  React.useEffect(() => {
    fetch('http://localhost:8000')
      .then((res) => {
        console.log(res);
        return res.json();
      })
      .then((res) => {
        console.log(res);
        setState(res.name);
      })
      .catch((e) => console.error(e));
  }, []);
  return (
    <div>
      <h1>Hello world</h1>
      <p>hello</p>
      <div>{state}</div>
    </div>
  );
}
