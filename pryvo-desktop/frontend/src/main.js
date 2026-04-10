import './style.css';
import './app.css';

import logo from './assets/images/pryvo__1_-removebg-preview.png';
import { GravarPerfil, VerificarPerfil, EscolherFoto } from '../wailsjs/go/main/App';


function desenharEcraRegisto() {
    document.querySelector('#app').innerHTML = `
        <div style="margin-top: 100px;">
            <img id="logo" style="width: 150px; margin-bottom: 20px;">
            <h1 style="color: #000000;">PRYVO</h1>
            <p style="color: black;">Welcome to your private social media. Create your local profile.</p>

            <div style="margin-top: 30px;">
                <input id="nomeInput" type="text" placeholder="What's your name?" style="padding: 10px; font-size: 16px; margin-bottom: 10px; width: 250px;" /><br>
                <div style="display: flex; align-items: center; justify-content: center; margin-bottom: 20px;">
                    <button onclick="processarFoto()" style="padding: 10px; font-size: 14px; cursor: pointer; margin-right: 15px; border-radius: 4px;">
                        Choose a picture
                    </button>

                    <img id="previewFoto" style="width: 45px; height: 45px; border-radius: 50%; display: none; object-fit: cover; border: 2px solid #000000;" />
                </div>

                <button onclick="criarPerfil()" style="padding: 10px 20px; font-size: 16px; cursor: pointer;">
                    Enter the system
                </button>
            </div>

            <div id="resultado" style="margin-top: 20px; color: red; font-weight: bold;"></div>
        </div>
    `;
    document.getElementById('logo').src = logo;
}


window.minhaFotoBase64 = "";

window.processarFoto = function() {
    EscolherFoto().then(imagemBase64 => {
        if (imagemBase64 !== "") {
            window.minhaFotoBase64 = imagemBase64; 
            let imgPreview = document.getElementById("previewFoto");
            imgPreview.src = window.minhaFotoBase64;
            imgPreview.style.display = "block";
        }
    }).catch(erro => {
        console.error("Erro a abrir janela nativa:", erro);
    });
};


function desenharMenuPrincipal(nomeDaPessoa, fotoDaPessoa) {
    let imagemSrc = fotoDaPessoa !== "" ? fotoDaPessoa : "https://cdn-icons-png.flaticon.com/512/149/149071.png";

    document.querySelector('#app').innerHTML = `
        <div style="margin-top: 80px;">
            <img src="${imagemSrc}" style="width: 100px; height: 100px; border-radius: 50%; object-fit: cover; border: 3px solid #000000; margin-bottom: 15px; box-shadow: 0 4px 8px rgba(0,0,0,0.2);">
            <h1 style="color: #000000; margin-top: 0;">Welcome back, ${nomeDaPessoa}!</h1>
            <p style="color: black;">Your profile is ready to use.</p>
            
            <div style="margin-top: 30px;">
                <button onclick="criarSalaNova('${nomeDaPessoa}', '${fotoDaPessoa}')" style="padding: 10px 20px; font-size: 16px; cursor: pointer; margin-right: 10px;">Create a room</button>
                <button onclick="entrarNaSala('${nomeDaPessoa}', '${fotoDaPessoa}')" style="padding: 10px 20px; font-size: 16px; cursor: pointer;">Enter a room</button>
            </div>
        </div>
    `;
}


window.criarPerfil = function () {
    let caixaNome = document.getElementById("nomeInput");
    let nomeLido = caixaNome.value.trim(); // O trim() corta os espaços vazios!
    let fotoLida = window.minhaFotoBase64;

    if (nomeLido === "") {
        document.getElementById("resultado").innerText = "Warning, you need to insert a name!";
        return;
    }

    GravarPerfil(nomeLido, fotoLida).then(sucesso => {
        desenharMenuPrincipal(nomeLido, fotoLida);
    }).catch(erro => {
        console.error("Erro de ligação ao Go:", erro);
    });
};


window.criarSalaNova = function(nomeDaPessoa, fotoDaPessoa) {
    let codigoAleatorio = Math.random().toString(36).substring(2, 7).toUpperCase();
    desenharEcraChat(codigoAleatorio, nomeDaPessoa, fotoDaPessoa);
};

window.entrarNaSala = function(nomeDaPessoa, fotoDaPessoa) {
    let codigoInserido = prompt("Enter the room code:");
    if (codigoInserido && codigoInserido.trim() !== "") {
        desenharEcraChat(codigoInserido.toUpperCase(), nomeDaPessoa, fotoDaPessoa);
    }
};

function desenharEcraChat(codigoSala, nomeDaPessoa, fotoDaPessoa) {
    document.querySelector('#app').innerHTML = `
        <div style="padding: 20px; max-width: 600px; margin: 0 auto; text-align: left;">
            <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #ccc; padding-bottom: 10px;">
                <h2 style="color: #ff3366; margin: 0;">Room: ${codigoSala}</h2>
                <button onclick="desenharMenuPrincipal('${nomeDaPessoa}', '${fotoDaPessoa}')" style="padding: 5px 15px; cursor: pointer;">Leave Room</button>
            </div>

            <div id="caixaMensagens" style="height: 350px; border: 1px solid #ccc; border-radius: 8px; margin-top: 20px; padding: 15px; overflow-y: auto; background: #ffffff;">
                <div style="color: #888; text-align: center; margin-top: 10px;">
                    Welcome to room ${codigoSala}.<br>Waiting for connection...
                </div>
            </div>

            <div style="display: flex; margin-top: 20px;">
                <input id="mensagemInput" type="text" placeholder="Type a message..." style="flex: 1; padding: 10px; font-size: 16px; border: 1px solid #ccc; border-radius: 4px;" />
                <button onclick="enviarMensagemLocal('${nomeDaPessoa}')" style="padding: 10px 20px; font-size: 16px; cursor: pointer; margin-left: 10px; background-color: #4CAF50; color: white; border: none; border-radius: 4px;">Send</button>
            </div>
        </div>
    `;

    document.getElementById("mensagemInput").addEventListener("keypress", function(event) {
        if (event.key === "Enter") {
            window.enviarMensagemLocal(nomeDaPessoa);
        }
    });
}

window.enviarMensagemLocal = function(nomeDaPessoa) {
    let caixaTexto = document.getElementById("mensagemInput");
    let textoEscrito = caixaTexto.value;
    if (textoEscrito.trim() === "") return;

    let ecraMensagens = document.getElementById("caixaMensagens");
    ecraMensagens.innerHTML += `
        <div style="margin-top: 15px; padding: 10px; background: #e3f2fd; border-radius: 8px; color: black; max-width: 80%;">
            <strong style="color: #1976d2;">${nomeDaPessoa}:</strong> ${textoEscrito}
        </div>
    `;
    caixaTexto.value = "";
    ecraMensagens.scrollTop = ecraMensagens.scrollHeight;
};


VerificarPerfil().then(perfil => {
    if (perfil.existe === true) {
        desenharMenuPrincipal(perfil.nome, perfil.foto);
    } else {
        desenharEcraRegisto();
    }
}).catch(erro => {
    console.error("Erro no Router:", erro);
    desenharEcraRegisto();
});