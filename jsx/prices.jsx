import React from 'react'
import Price from './price'

export default class Prices extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            data: []
            new: {}
        };

        if ('data' in props) {
            this.state.data = props.data;
        }
    }

    render() {
        return (
            <React.Fragment>
                this.state.data.map((price) =>
                    <Price key={price.ID} data={price} />
                )
            </React.Fragment>
        )
    }
}
