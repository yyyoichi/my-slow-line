import 'styles/globals.css';
import dynamic from 'next/dynamic';

export default dynamic(() => import('templates/_app/'), { ssr: false });
