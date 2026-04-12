import './style.css';
import './app.css';

import logo from './assets/images/pryvo__1_-removebg-preview.png';
import { GravarPerfil, VerificarPerfil, EscolherFoto, CriarSalaNet, EntrarSalaNet, EnviarMensagemNet, SairDaSalaNet } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';

function desenharEcraRegisto() {
    document.querySelector('#app').innerHTML = `
        <div class="app-card">
            <img id="logo" src="${logo}" class="logo-img">
            <h1>PRYVO</h1>
            <p>Welcome to your private social media. Create your local profile.</p>

            <input id="nomeInput" type="text" placeholder="What's your name?" class="input-modern" />
            
            <div style="display: flex; align-items: center; justify-content: center; margin-bottom: 20px;">
                <button onclick="processarFoto()" class="btn-secondary">
                    📷 Choose a picture
                </button>
                <img id="previewFoto" class="avatar-preview" style="display: none;" />
            </div>

            <button onclick="criarPerfil()" class="btn-primary">
                Enter the system
            </button>

            <div id="resultado" style="margin-top: 15px; color: var(--danger); font-weight: bold; font-size: 14px;"></div>
        </div>
    `;
}

window.minhaFotoBase64 = "";

window.processarFoto = function() {
    EscolherFoto().then(imagemBase64 => {
        if (imagemBase64 !== "") {
            let img = new Image();
            img.onload = function() {
                let canvas = document.createElement('canvas');
                let tamanho = 100;
                canvas.width = tamanho;
                canvas.height = tamanho;
                let ctx = canvas.getContext('2d');
                ctx.drawImage(img, 0, 0, tamanho, tamanho);
                window.minhaFotoBase64 = canvas.toDataURL('image/jpeg', 0.5);

                let imgPreview = document.getElementById("previewFoto");
                imgPreview.src = window.minhaFotoBase64;
                imgPreview.style.display = "block";
            };
            img.src = imagemBase64;
        }
    }).catch(erro => console.error("Erro a abrir janela nativa:", erro));
};

window.desenharMenuPrincipal = function(nomeDaPessoa, fotoDaPessoa) {
    //let imagemSrc = fotoDaPessoa !== "" ? fotoDaPessoa : "https://cdn-icons-png.flaticon.com/512/149/149071.png";
    let imagemSrc = (fotoDaPessoa && fotoDaPessoa.trim() !== "") ? fotoDaPessoa : "https://cdn-icons-png.flaticon.com/512/149/149071.png";

    document.querySelector('#app').innerHTML = `
        <div class="app-card">
            <img src="${imagemSrc}" class="avatar-large">
            <h1>Welcome back, ${nomeDaPessoa}!</h1>
            <p>Your profile is ready to use.</p>
            
            <div style="display: flex; gap: 10px; margin-top: 20px;">
                <button onclick="criarSalaNova('${nomeDaPessoa}', '${fotoDaPessoa}')" class="btn-primary" style="flex: 1;">Create Room</button>
                <button onclick="entrarNaSala('${nomeDaPessoa}', '${fotoDaPessoa}')" class="btn-secondary" style="flex: 1;">Join Room</button>
            </div>
        </div>
    `;
}

window.criarPerfil = function () {
    let caixaNome = document.getElementById("nomeInput");
    let nomeLido = caixaNome.value.trim(); 
    let fotoLida = window.minhaFotoBase64;

    if (nomeLido === "") {
        document.getElementById("resultado").innerText = "Warning: You need to insert a name!";
        return;
    }

    GravarPerfil(nomeLido, fotoLida).then(sucesso => {
        window.desenharMenuPrincipal(nomeLido, fotoLida);
    }).catch(erro => console.error("Erro de ligação ao Go:", erro));
};

window.criarSalaNova = function(nomeDaPessoa, fotoDaPessoa) {
    let codigoAleatorio = Math.random().toString(36).substring(2, 7).toUpperCase();
    CriarSalaNet(codigoAleatorio, nomeDaPessoa, fotoDaPessoa).then(respostaDoGo => {
        if (respostaDoGo === "ok") {
            window.desenharEcraChat(codigoAleatorio, nomeDaPessoa, fotoDaPessoa);
        } else {
            alert("Atenção: " + respostaDoGo); 
        }
    }).catch(erro => console.error(erro));
};

window.entrarNaSala = function(nomeDaPessoa, fotoDaPessoa) {
    let codigoInserido = prompt("Enter the room code:");
    if (codigoInserido && codigoInserido.trim() !== "") {
        codigoInserido = codigoInserido.toUpperCase();
        EntrarSalaNet(codigoInserido, nomeDaPessoa, fotoDaPessoa).then(respostaDoGo => {
            if (respostaDoGo === "ok") {
                window.desenharEcraChat(codigoInserido, nomeDaPessoa, fotoDaPessoa);
            } else {
                alert("Não foi possível entrar: " + respostaDoGo);
            }
        }).catch(erro => console.error(erro));
    }
};

window.sairDaSala = function(nomeDaPessoa, fotoDaPessoa) {
    SairDaSalaNet().then(() => {
        window.desenharMenuPrincipal(nomeDaPessoa, fotoDaPessoa);
    });
};

window.desenharEcraChat = function(codigoSala, nomeDaPessoa, fotoDaPessoa) {
    document.querySelector('#app').innerHTML = `
        <div class="chat-wrapper">
            <header class="chat-header">
                <div class="chat-header-info">
                    <img src="${fotoDaPessoa !== '' ? fotoDaPessoa : 'https://cdn-icons-png.flaticon.com/512/149/149071.png'}" class="avatar-small">
                    <h2>Room: ${codigoSala}</h2>
                </div>
                <button onclick="sairDaSala('${nomeDaPessoa}', '${fotoDaPessoa}')" class="btn-danger">Leave</button>
            </header>

            <div id="caixaMensagens" class="chat-messages">
                <div class="system-message" id="avisoEspera">
                    Welcome to room ${codigoSala}.<br>Waiting for connection...
                </div>
            </div>

            <footer class="chat-footer">
                <input id="mensagemInput" type="text" placeholder="Type a message..." class="chat-input" />
                <button onclick="enviarMensagemLocal('${nomeDaPessoa}', '${fotoDaPessoa}')" class="btn-send">
                    <svg viewBox="0 0 24 24"><path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"></path></svg>
                </button>
            </footer>
        </div>
    `;

    document.getElementById("mensagemInput").addEventListener("keypress", function(event) {
        if (event.key === "Enter") {
            window.enviarMensagemLocal(nomeDaPessoa, fotoDaPessoa);
        }
    });
}

window.enviarMensagemLocal = function(nomeDaPessoa, fotoDaPessoa) {
    let caixaTexto = document.getElementById("mensagemInput");
    let textoEscrito = caixaTexto.value;
    if (textoEscrito.trim() === "") return;

    EnviarMensagemNet(textoEscrito);

    let ecraMensagens = document.getElementById("caixaMensagens");
    ecraMensagens.innerHTML += `
        <div class="msg-row mine">
            <div class="msg-bubble">
                <div class="msg-text">${textoEscrito}</div>
            </div>
        </div>
    `;
    
    caixaTexto.value = "";
    ecraMensagens.scrollTop = ecraMensagens.scrollHeight;
};

EventsOn("room_ready", function(mensagemDoGo) {
    let aviso = document.getElementById("avisoEspera");
    if (aviso) {
        aviso.innerHTML = `<strong>${mensagemDoGo}</strong>`;
        aviso.style.color = "#00a884";
        aviso.style.backgroundColor = "#e6fce5";
    }
});

EventsOn("mensagem_recebida", function(msg) {
    let ecraMensagens = document.getElementById("caixaMensagens");
    if (!ecraMensagens) return;

    ecraMensagens.innerHTML += `
        <div class="msg-row theirs">
            <div class="msg-bubble">
                <div class="msg-header">
                    <img src="${msg.foto !== "" ? msg.foto : 'https://cdn-icons-png.flaticon.com/512/149/149071.png'}" class="msg-avatar">
                    ${msg.nome}
                </div>
                <div class="msg-text">${msg.texto}</div>
            </div>
        </div>
    `;
    ecraMensagens.scrollTop = ecraMensagens.scrollHeight;
});

VerificarPerfil().then(perfil => {
    if (perfil.existe === true) {
        window.desenharMenuPrincipal(perfil.nome, perfil.foto);
    } else {
        desenharEcraRegisto();
    }
}).catch(erro => {
    console.error("Erro no Router:", erro);
    desenharEcraRegisto();
});