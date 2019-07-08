import { Component } from 'react';
import StateErrors from './stateErrors';

interface ValidateFunc {
	(name: string, value: string): boolean;
}

interface ValueFunc {
	(name: string, value: string): Object;
}

export default class FormHelper extends StateErrors {
	constructor(props) {
		super(props);
		
		if(!this.state) {
			this.state = {}
		}
		
		if(!this.state.hasOwnProperty("form")) {
			this.state.form = {
				values: {},
				validations: {},
			};
		}
		
		this.validate = this.validate.bind(this);
		this.onChange = this.onChange.bind(this);
		this.getValue = this.getValue.bind(this);
		this.getValues = this.getValues.bind(this);
		this.mapValues = this.mapValues.bind(this);
	}

	validate(name: RegExp, fn: ValidateFunc) {
		this.setState((state) => {
			if(!state.form.validations.hasOwnProperty(name)) {
				state.form.validations[name] = []; 
			}
			
			state.form.validations[name].push(fn);
			return state;
		} );
	}
	
	getValue(name: string) {
		if(this.state.form.values.hasOwnProperty(name)) {
			return this.state.form.values[name];
		}
	}
	
	getValues() {
		return this.state.form.values;
	}
	
	mapValues(fn: ValueFunc) {
		return Object.keys(this.state.form.values).map((name) => {
			return fn(name, this.state.form.values[name]);
		});
	}
	
	onChange(e, { name, value }) {
		e.preventDefault();
		
		if(Object.keys(this.state.form.validations).map((regex) => {
			if(name.match(regex)) return ! this.state.form.validations[regex].map(
				(fn) => fn(name, value)
			).includes(false);
		}).includes(false)) return;
		
		this.onError(name, false);
	
		this.setState((state) => {
			state.form.values[name] = value;
			return state;
		} );
	}
}
