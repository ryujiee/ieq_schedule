<template>
    <q-page class="q-pa-md">

        <!-- Header -->
        <div class="row items-center q-col-gutter-md q-mb-md">
            <div class="col-12 col-md">
                <div class="text-h5 text-weight-bold">WhatsApp</div>
                <div class="text-caption text-grey-7">
                    Conecte o WhatsApp via QR Code para enviar lembretes automáticos
                </div>
            </div>

            <div class="col-12 col-md-auto">
                <div class="row items-center justify-end q-gutter-sm">
                    <q-btn icon="refresh" label="Atualizar" unelevated no-caps @click="refreshAll" :loading="loading" />
                    <q-btn color="primary" icon="qr_code_2" label="Conectar" unelevated no-caps @click="startConnect"
                        :disable="status.loggedIn && status.connected" :loading="connecting" />
                    <q-btn color="negative" icon="logout" label="Logout" unelevated no-caps @click="confirmLogout"
                        :disable="!status.loggedIn" />
                </div>
            </div>
        </div>

        <div class="row q-col-gutter-md">

            <!-- Status -->
            <div class="col-12 col-md-5">
                <q-card bordered class="panel-card">
                    <q-card-section class="row items-center">
                        <div class="text-subtitle1 text-weight-bold">Status</div>
                        <q-space />
                        <q-badge outline :color="status.connected ? 'positive' : 'grey-7'">
                            {{ status.connected ? 'Conectado' : 'Desconectado' }}
                        </q-badge>
                    </q-card-section>

                    <q-separator />

                    <q-card-section>
                        <div class="q-gutter-sm">
                            <div class="row items-center no-wrap">
                                <q-icon name="verified_user" class="q-mr-sm text-grey-7" />
                                <div class="text-weight-medium">Sessão</div>
                                <q-space />
                                <q-badge outline :color="status.loggedIn ? 'positive' : 'grey-7'">
                                    {{ status.loggedIn ? 'Ativa' : 'Não conectada' }}
                                </q-badge>
                            </div>

                            <div class="row items-center no-wrap">
                                <q-icon name="sync" class="q-mr-sm text-grey-7" />
                                <div class="text-weight-medium">Pairing (QR)</div>
                                <q-space />
                                <q-badge outline :color="status.pairing ? 'primary' : 'grey-7'">
                                    {{ status.pairing ? 'Em andamento' : 'Parado' }}
                                </q-badge>
                            </div>

                            <q-banner rounded class="bg-grey-2 text-grey-9 q-mt-md">
                                <template #avatar>
                                    <q-icon name="info" color="primary" />
                                </template>

                                <div v-if="status.loggedIn && status.connected">
                                    ✅ WhatsApp conectado! Agora o sistema poderá enviar lembretes automáticos.
                                </div>

                                <div v-else-if="status.pairing">
                                    📷 Escaneie o QR Code no seu WhatsApp:
                                    <b>Menu → Dispositivos conectados → Conectar um dispositivo</b>
                                </div>

                                <div v-else>
                                    Clique em <b>Conectar</b> para gerar um QR Code e parear o WhatsApp.
                                </div>
                            </q-banner>
                        </div>
                    </q-card-section>
                </q-card>

                <!-- Dica -->
                <div class="text-caption text-grey-7 q-mt-sm">
                    Dica: use um número oficial da mídia (evite o pessoal) para não misturar conversas.
                </div>
            </div>

            <!-- QR -->
            <div class="col-12 col-md-7">
                <q-card bordered class="panel-card">
                    <q-card-section class="row items-center">
                        <div class="text-subtitle1 text-weight-bold">QR Code</div>
                        <q-space />
                        <q-btn flat dense icon="content_copy" label="Copiar código" no-caps @click="copyRawQR"
                            :disable="!qrRaw" />
                    </q-card-section>

                    <q-separator />

                    <q-card-section class="q-pa-lg">

                        <div v-if="loadingQr" class="row items-center q-gutter-sm">
                            <q-spinner size="28px" color="primary" />
                            <div class="text-grey-7">Carregando QR...</div>
                        </div>

                        <div v-else-if="status.loggedIn && status.connected" class="column items-center text-center">
                            <q-icon name="check_circle" size="64px" class="text-positive q-mb-sm" />
                            <div class="text-subtitle1 text-weight-bold">Tudo certo!</div>
                            <div class="text-caption text-grey-7">
                                WhatsApp já está conectado. Não é necessário QR agora.
                            </div>
                        </div>

                        <div v-else-if="qrDataUrl" class="column items-center">
                            <img :src="qrDataUrl" alt="QR Code WhatsApp" class="qr-img" />
                            <div class="text-caption text-grey-7 q-mt-md text-center">
                                Abra o WhatsApp no celular e escaneie este QR Code.
                            </div>
                        </div>

                        <div v-else class="column items-center text-center text-grey-7">
                            <q-icon name="qr_code_2" size="64px" class="q-mb-sm" />
                            <div class="text-subtitle2 text-weight-bold">Sem QR disponível</div>
                            <div class="text-caption q-mt-xs">
                                Clique em <b>Conectar</b> para gerar um QR Code.
                            </div>
                        </div>

                    </q-card-section>
                </q-card>
            </div>

        </div>
    </q-page>
</template>

<script>
import { api } from 'boot/api'
import QRCode from 'qrcode'

export default {
    name: 'AdminWhatsApp',

    data() {
        return {
            loading: false,
            connecting: false,
            loadingQr: false,

            status: {
                connected: false,
                loggedIn: false,
                pairing: false,
                lastQr: ''
            },

            qrRaw: '',
            qrDataUrl: '',

            pollTimer: null
        }
    },

    mounted() {
        this.refreshAll()
        this.startPolling()
    },

    beforeUnmount() {
        this.stopPolling()
    },

    methods: {
        startPolling() {
            this.stopPolling()
            this.pollTimer = setInterval(() => {
                // se já está conectado/logado, podemos reduzir chamadas
                this.fetchStatus()
                if (!this.status.loggedIn || !this.status.connected) {
                    this.fetchQR()
                }
            }, 1200)
        },

        stopPolling() {
            if (this.pollTimer) {
                clearInterval(this.pollTimer)
                this.pollTimer = null
            }
        },

        async refreshAll() {
            this.loading = true
            try {
                await this.fetchStatus(true)
                await this.fetchQR(true)
            } finally {
                this.loading = false
            }
        },

        async fetchStatus(silent = false) {
            try {
                const { data } = await api.get('/admin/whatsapp/status')
                this.status = {
                    connected: !!data.connected,
                    loggedIn: !!data.loggedIn,
                    pairing: !!data.pairing,
                    lastQr: data.lastQr || ''
                }
            } catch (e) {
                console.error(e)
                if (!silent) this.$q.notify({ type: 'negative', message: 'Erro ao buscar status do WhatsApp' })
            }
        },

        async fetchQR(silent = false) {
            this.loadingQr = true
            try {
                const { data } = await api.get('/admin/whatsapp/qr')

                const qr = data.qr || ''
                this.qrRaw = qr

                if (qr) {
                    // só regen se mudou
                    if (!this.qrDataUrl || qr !== this.status.lastQr) {
                        this.qrDataUrl = await QRCode.toDataURL(qr, { margin: 1, scale: 7 })
                    }
                } else {
                    this.qrDataUrl = ''
                }
            } catch (e) {
                console.error(e)

                if (!silent) this.$q.notify({ type: 'negative', message: 'Erro ao buscar QR do WhatsApp' })
            } finally {
                this.loadingQr = false
            }
        },

        async startConnect() {
            this.connecting = true
            try {
                await api.post('/admin/whatsapp/connect')
                await this.fetchStatus(true)
                await this.fetchQR(true)
                this.$q.notify({ type: 'positive', message: 'Pairing iniciado. Escaneie o QR Code.' })
            } catch (e) {
                console.error(e)

                this.$q.notify({ type: 'negative', message: 'Falha ao iniciar conexão do WhatsApp' })
            } finally {
                this.connecting = false
            }
        },

        copyRawQR() {
            if (!this.qrRaw) return
            navigator.clipboard.writeText(this.qrRaw)
            this.$q.notify({ type: 'positive', message: 'Código copiado' })
        },

        confirmLogout() {
            this.$q.dialog({
                title: 'Logout do WhatsApp',
                message: 'Isso vai remover a sessão salva. Para usar novamente será necessário escanear o QR Code. Continuar?',
                cancel: true,
                persistent: true
            }).onOk(() => this.logout())
        },

        async logout() {
            try {
                await api.post('/admin/whatsapp/logout')
                this.qrRaw = ''
                this.qrDataUrl = ''
                await this.fetchStatus(true)
                this.$q.notify({ type: 'positive', message: 'Sessão removida. Conecte novamente via QR.' })
            } catch (e) {
                console.error(e)

                this.$q.notify({ type: 'negative', message: 'Erro ao fazer logout do WhatsApp' })
            }
        }
    }
}
</script>

<style scoped>
.panel-card {
    border-radius: 14px;
}

.qr-img {
    width: 320px;
    max-width: 92vw;
    border-radius: 14px;
    border: 1px solid rgba(0, 0, 0, .08);
    padding: 12px;
    background: white;
}
</style>
