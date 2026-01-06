<template>
  <q-page class="q-pa-md">

    <!-- Header -->
    <div class="row items-center q-col-gutter-md q-mb-md">
      <div class="col-12 col-md">
        <div class="text-h5 text-weight-bold">Configurações</div>
        <div class="text-caption text-grey-7">
          Ajuste o horário e a mensagem dos lembretes automáticos
        </div>
      </div>

      <div class="col-12 col-md-auto">
        <div class="row items-center justify-end q-gutter-sm">
          <q-btn
            flat
            icon="refresh"
            label="Recarregar"
            no-caps
            @click="load"
            :loading="loading"
          />
          <q-btn
            color="primary"
            icon="save"
            label="Salvar"
            unelevated
            no-caps
            @click="save"
            :loading="saving"
          />
        </div>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Painel principal -->
      <div class="col-12 col-md-7">
        <q-card bordered class="panel-card">
          <q-card-section class="row items-center">
            <div class="text-subtitle1 text-weight-bold">Lembretes</div>
            <q-space />
            <q-toggle v-model="form.remindersEnabled" label="Ativar" />
          </q-card-section>

          <q-separator />

          <q-card-section class="q-gutter-md">
            <div class="row q-col-gutter-md">
              <div class="col-12 col-sm-6">
                <q-input
                  v-model.number="form.reminderHour"
                  type="number"
                  outlined
                  dense
                  label="Hora do lembrete (0-23)"
                  :disable="!form.remindersEnabled"
                  :rules="[
                    v => v !== null && v !== '' || 'Obrigatório',
                    v => v >= 0 && v <= 23 || 'Use 0 a 23'
                  ]"
                >
                  <template #prepend><q-icon name="schedule" /></template>
                </q-input>
              </div>

              <div class="col-12 col-sm-6">
                <q-input
                  v-model.number="form.reminderMinute"
                  type="number"
                  outlined
                  dense
                  label="Minuto (0-59)"
                  :disable="!form.remindersEnabled"
                  :rules="[
                    v => v !== null && v !== '' || 'Obrigatório',
                    v => v >= 0 && v <= 59 || 'Use 0 a 59'
                  ]"
                >
                  <template #prepend><q-icon name="timer" /></template>
                </q-input>
              </div>
            </div>

            <q-banner rounded class="bg-grey-2 text-grey-9">
              <template #avatar>
                <q-icon name="info" color="primary" />
              </template>
              O horário sempre será interpretado no fuso de <b>Brasília (America/Sao_Paulo)</b>.
            </q-banner>
          </q-card-section>
        </q-card>

        <q-card bordered class="panel-card q-mt-md">
          <q-card-section class="row items-center">
            <div class="text-subtitle1 text-weight-bold">Mensagem do WhatsApp</div>
            <q-space />
            <q-badge outline color="primary">variáveis</q-badge>
          </q-card-section>

          <q-separator />

          <q-card-section class="q-gutter-md">
            <q-input
              v-model="form.reminderMessage"
              type="textarea"
              outlined
              autogrow
              label="Mensagem padrão"
              hint="Use {Nome} e {Funcao}"
              :disable="!form.remindersEnabled"
            >
              <template #prepend><q-icon name="chat" /></template>
            </q-input>

            <div class="row q-gutter-sm">
              <q-btn
                flat
                no-caps
                icon="person"
                label="Inserir {Nome}"
                @click="insertToken('{Nome}')"
                :disable="!form.remindersEnabled"
              />
              <q-btn
                flat
                no-caps
                icon="badge"
                label="Inserir {Funcao}"
                @click="insertToken('{Funcao}')"
                :disable="!form.remindersEnabled"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Preview -->
      <div class="col-12 col-md-5">
        <q-card bordered class="panel-card">
          <q-card-section class="row items-center">
            <div class="text-subtitle1 text-weight-bold">Prévia</div>
            <q-space />
            <q-badge outline color="grey-8">exemplo</q-badge>
          </q-card-section>

          <q-separator />

          <q-card-section>
            <div class="text-caption text-grey-7 q-mb-sm">
              Exemplo de como a pessoa receberá:
            </div>

            <div class="wa-preview">
              <div class="wa-bubble">
                {{ previewMessage }}
              </div>
              <div class="text-caption text-grey-6 q-mt-xs">
                {{ previewTime }}
              </div>
            </div>

            <q-separator class="q-my-md" />

            <div class="text-caption text-grey-7">
              Dica: deixe a mensagem curta e clara.
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

  </q-page>
</template>

<script>
import { api } from 'boot/api'

export default {
  name: 'AdminSettings',

  data () {
    return {
      loading: false,
      saving: false,
      form: {
        remindersEnabled: true,
        reminderHour: 8,
        reminderMinute: 0,
        reminderMessage: 'Olá {Nome}! Hoje é o seu dia de servir na função {Funcao}. Deus abençoe! 🙏'
      }
    }
  },

  computed: {
    previewMessage () {
      const msg = this.form.reminderMessage || ''
      return msg
        .replaceAll('{Nome}', 'João')
        .replaceAll('{Funcao}', 'Transmissão')
        .trim() || '—'
    },

    previewTime () {
      const h = String(this.form.reminderHour ?? 8).padStart(2, '0')
      const m = String(this.form.reminderMinute ?? 0).padStart(2, '0')
      return `Hoje • ${h}:${m}`
    }
  },

  mounted () {
    this.load()
  },

  methods: {
    async load () {
      this.loading = true
      try {
        const { data } = await api.get('/admin/settings')
        // espera algo tipo:
        // { remindersEnabled, reminderHour, reminderMinute, reminderMessage }
        this.form = {
          remindersEnabled: !!data.remindersEnabled,
          reminderHour: Number(data.reminderHour ?? 8),
          reminderMinute: Number(data.reminderMinute ?? 0),
          reminderMessage: data.reminderMessage || this.form.reminderMessage
        }
      } catch (e) {
        console.error(e)
        this.$q.notify({ type: 'negative', message: 'Falha ao carregar configurações' })
      } finally {
        this.loading = false
      }
    },

    async save () {
      // valida simples
      if (this.form.reminderHour < 0 || this.form.reminderHour > 23) {
        this.$q.notify({ type: 'warning', message: 'Hora inválida (0-23)' })
        return
      }
      if (this.form.reminderMinute < 0 || this.form.reminderMinute > 59) {
        this.$q.notify({ type: 'warning', message: 'Minuto inválido (0-59)' })
        return
      }

      this.saving = true
      try {
        await api.put('/admin/settings', {
          remindersEnabled: !!this.form.remindersEnabled,
          reminderHour: Number(this.form.reminderHour),
          reminderMinute: Number(this.form.reminderMinute),
          reminderMessage: (this.form.reminderMessage || '').trim()
        })

        this.$q.notify({ type: 'positive', message: 'Configurações salvas!' })
      } catch (e) {
        console.error(e)
        this.$q.notify({ type: 'negative', message: 'Erro ao salvar configurações' })
      } finally {
        this.saving = false
      }
    },

    insertToken (token) {
      this.form.reminderMessage = (this.form.reminderMessage || '') + token
    }
  }
}
</script>

<style scoped>
.panel-card {
  border-radius: 14px;
}

.wa-preview {
  background: #f7f7fb;
  border: 1px solid rgba(0,0,0,.06);
  border-radius: 14px;
  padding: 14px;
}

.wa-bubble {
  background: white;
  border: 1px solid rgba(0,0,0,.08);
  border-radius: 14px;
  padding: 12px 12px;
  max-width: 100%;
  white-space: pre-line;
}
</style>
