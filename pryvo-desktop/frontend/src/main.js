import './style.css';
import './app.css';

import logo from './assets/images/pryvo__1_-removebg-preview.png';
import { GravarPerfil, VerificarPerfil } from '../wailsjs/go/main/App';

function desenharEcraRegisto() {
    document.querySelector('#app').innerHTML = `
        <div style="margin-top: 100px;">
            <img id="logo" style="width: 150px; margin-bottom: 20px;">
            <h1 style="color: #000000;">PRYVO</h1>
            <p style="color: black;">Welcome to your private social media. Create your local profile.</p>

            <div style="margin-top: 30px;">
                <input id="nomeInput" type="text" placeholder="What's your name?" style="padding: 10px; font-size: 16px; margin-bottom: 10px; width: 250px;" /><br>
                <input id="fotoInput" type="text" placeholder="Choose a picture(optional)..." style="padding: 10px; font-size: 16px; margin-bottom: 10px; width: 250px;" /><br>

                <button onclick="criarPerfil()" style="padding: 10px 20px; font-size: 16px; cursor: pointer;">
                    Enter the system
                </button>
            </div>

            <div id="resultado" style="margin-top: 20px; color: #4caf50; font-weight: bold;"></div>
        </div>
    `;
    document.getElementById('logo').src = logo;
}

function desenharMenuPrincipal(nomeDaPessoa) {
    document.querySelector('#app').innerHTML = `
        <div style="margin-top: 100px;">
            <h1 style="color: #4CAF50;">Welcome back, ${nomeDaPessoa}!</h1>
            <p style="color: black;">Your profile is ready to use.</p>
            
            <div style="margin-top: 30px;">
                <button style="padding: 10px 20px; font-size: 16px; cursor: pointer; margin-right: 10px;">Create a room</button>
                <button style="padding: 10px 20px; font-size: 16px; cursor: pointer;">Enter a room</button>
            </div>
        </div>
    `;
}

window.criarPerfil = function () {
    let nomeLido = document.getElementById("nomeInput").value;
    let fotoLida = document.getElementById("fotoInput").value;

    if (nomeLido === "") {
        document.getElementById("resultado").innerText = "Warning, you need to insert a name!";
        document.getElementById("resultado").style.color = "red";
        return;
    }

    GravarPerfil(nomeLido, fotoLida).then(sucesso => {
        desenharMenuPrincipal(nomeLido);
    }).catch(erro => {
        console.error("Erro de ligação ao Go:", erro);
    });
};

VerificarPerfil().then(perfil => {
    if (perfil.existe === true) {
        desenharMenuPrincipal(perfil.nome);
    } else {
        desenharEcraRegisto();
    }
}).catch(erro => {
    console.error("Erro no Router:", erro);
    desenharEcraRegisto();
});