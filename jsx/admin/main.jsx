import React from 'react'
import ReactDOM from 'react-dom'
import Items from './items'

import '../../semantic/dist/semantic.min.css';

ReactDOM.render(
    <React.Fragment>
        <h1>Hello, world!</h1><br />
        <Items />
    </React.Fragment>,
    document.getElementById('root')
);
