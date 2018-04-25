// src/index.ts

import Vue from "vue";
import HelloComponent from './components/hello.vue';
import * as firebase from 'firebase';

var config = {
    apiKey: "AIzaSyCn835yngATs0mTLB3aApJI2d-TMRalap4",
    authDomain: "ping-party.firebaseapp.com",
    projectId: "ping-party"
};
firebase.initializeApp(config);

var app: Vue;
firebase.auth().onAuthStateChanged((x) => {
    if (!app) {
        app = new Vue({
            el: "#app",
            template: `
            <div>
                <div>Hello {{name}}!</div>
                Name: <input v-model="name" type="text">
                <hello-component :name="name" :initialEnthusiasm="5"/>
            </div>
            `,
            data: {
                name: "World"
            },
            components: {
                HelloComponent
            }
        });
    }
})
