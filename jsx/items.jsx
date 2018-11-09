import React from 'react'
import Item from './item'

export default class Items extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            data: []
        };

        if ('data' in props) {
            this.state.data = props.data;
        } else {
            this.fetchItems();
        }
    }

    fetchItems() {
        fetch("/shop/keeper/badges.json", {
            method: "GET"
        }).then((response) => response.json()).then((data) => {
            this.setState((state) => {
                state.data = data
                return state
            })
        })
    }

    render() {
        return this.state.data.map((item) =>
            <Item key={item.ID} data={item} />
        )
    }
}
