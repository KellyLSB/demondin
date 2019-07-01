import { Component } from 'react';

export default class StateErrors extends Component {
	constructor(props) {
		super(props);
		
		if(!this.state) {
			this.state = {}
		}
		
		if(!this.state.hasOwnProperty("errors")) {
			this.state.errors = {};
		}
		
		this.onError = this.onError.bind(this);
		this.hasError = this.hasError.bind(this);
		this.getErrors = this.getErrors.bind(this);
	}

	onError(name, msg = null) {
		this.setState((state) => {
			if (msg === false) {
				delete state.errors[name];
			} else {
				state.errors[name] = msg;
			}
			
			return state;
		} );
	}
	
	hasError(name = null) {
		if (name === null) {
			return Object.keys(this.state.errors).length > 0;
		}
		
		return this.state.errors.hasOwnProperty(name);
	}
	
	getErrors() {
		return Object.values(this.state.errors);
	}
}
