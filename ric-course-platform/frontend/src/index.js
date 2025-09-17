import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import './styles.css';
import './background.css';


const container = document.getElementById('root');
const root = createRoot(container);
root.render(<App />);
