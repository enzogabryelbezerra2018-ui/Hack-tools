package com.simuladorbackup

import android.app.Activity
import android.app.AlertDialog
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.provider.DocumentsContract
import android.util.Log
import android.widget.Button
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import androidx.documentfile.provider.DocumentFile
import kotlinx.coroutines.*
import java.io.BufferedInputStream
import java.io.BufferedOutputStream
import java.io.File
import java.io.FileOutputStream
import java.io.InputStream

class MainActivity : AppCompatActivity() {
    private lateinit var btnStart: Button
    private lateinit var logsView: TextView
    private val REQUEST_CODE_PICK_DIR = 1001
    private val ioScope = CoroutineScope(Dispatchers.IO + SupervisorJob())

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main) // layout simples com btnStart e logsView

        btnStart = findViewById(R.id.btnStart)
        logsView = findViewById(R.id.logs)

        btnStart.setOnClickListener {
            showConfirmationDialog()
        }

        // Se Activity iniciada por UsbConnectedReceiver, mostramos diálogo automático
        val triggered = intent.getBooleanExtra("triggered_by_usb", false)
        if (triggered) {
            showConfirmationDialog()
        }
    }

    private fun showConfirmationDialog() {
        runOnUiThread {
            AlertDialog.Builder(this)
                .setTitle("Atenção")
                .setMessage("Você está com seu celular de teste?\n(Escolha 'Não' para cancelar)")
                .setNegativeButton("Não") { dlg, _ ->
                    appendLog("Usuário cancelou (Não).", "RED")
                    dlg.dismiss()
                }
                .setPositiveButton("Sim") { dlg, _ ->
                    dlg.dismiss()
                    appendLog("Usuário aceitou (Sim). Abra o seletor e escolha a raiz que deseja copiar.", "YELLOW")
                    openDirectoryPicker()
                }
                .setCancelable(false)
                .show()
        }
    }

    private fun openDirectoryPicker() {
        // ACTION_OPEN_DOCUMENT_TREE -> usuário escolhe a pasta/raiz (ex: "Internal shared storage")
        val intent = Intent(Intent.ACTION_OPEN_DOCUMENT_TREE)
        // opcional: mostrar opções para o usuário; deixamos padrão
        startActivityForResult(intent, REQUEST_CODE_PICK_DIR)
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        if (requestCode == REQUEST_CODE_PICK_DIR && resultCode == Activity.RESULT_OK) {
            val treeUri: Uri? = data?.data
            if (treeUri == null) {
                appendLog("Nenhuma pasta selecionada.", "RED")
                return
            }
            // Persistir permissão para ler a pasta escolhida (recomendado)
            val takeFlags: Int = (Intent.FLAG_GRANT_READ_URI_PERMISSION
                    or Intent.FLAG_GRANT_WRITE_URI_PERMISSION)
            contentResolver.takePersistableUriPermission(treeUri, takeFlags)

            appendLog("Pasta selecionada: $treeUri", "YELLOW")
            ioScope.launch {
                performCopyToDataFolder(treeUri)
            }
        } else {
            super.onActivityResult(requestCode, resultCode, data)
        }
    }

    private fun appendLog(text: String, level: String = "NORMAL") {
        runOnUiThread {
            val prefix = when(level) {
                "RED" -> "⛔ $text"
                "YELLOW" -> "… $text"
                "GREEN" -> "✔ $text"
                else -> text
            }
            logsView.append(prefix + "\n")
        }
    }

    /**
     * Copia tudo que for acessível via DocumentFile (SAF) -> para getExternalFilesDir(null)/data
     * Se arquivos com mesmo nome já existirem, adiciona sufixo _1, _2, ... para evitar sobrescrita.
     */
    private suspend fun performCopyToDataFolder(treeUri: Uri) {
        withContext(Dispatchers.IO) {
            try {
                appendLog("Iniciando cópia para pasta 'data' do app...", "YELLOW")
                val sourceTree = DocumentFile.fromTreeUri(this@MainActivity, treeUri)
                if (sourceTree == null || !sourceTree.isDirectory) {
                    appendLog("Seleção inválida (não é diretório).", "RED")
                    return@withContext
                }

                val destBase: File = getExternalFilesDir(null) ?: filesDir
                val destDataFolder = File(destBase, "data")
                if (!destDataFolder.exists()) destDataFolder.mkdirs()
                appendLog("Destino (local): ${destDataFolder.absolutePath}", "YELLOW")

                copyDocumentFileTreeToFile(sourceTree, destDataFolder)

                appendLog("Cópia finalizada.", "GREEN")
                appendLog("OBS: Arquivos originais não foram apagados nem enviados.", "GREEN")
            } catch (e: Exception) {
                appendLog("Erro: ${e.message}", "RED")
            }
        }
    }

    /**
     * Percorre DocumentFile source e copia recursivamente para dest (File).
     */
    private fun copyDocumentFileTreeToFile(source: DocumentFile, dest: File) {
        if (!dest.exists()) dest.mkdirs()
        val children = source.listFiles()
        for (child in children) {
            // small safety check: ignore null names
            val childName = child.name ?: continue
            if (child.isDirectory) {
                val subdir = File(dest, childName)
                subdir.mkdirs()
                copyDocumentFileTreeToFile(child, subdir)
            } else if (child.isFile) {
                val outFile = uniqueFile(File(dest, childName))
                try {
                    contentResolver.openInputStream(child.uri)?.use { input: InputStream ->
                        BufferedInputStream(input).use { bis ->
                            FileOutputStream(outFile).use { fos ->
                                BufferedOutputStream(fos).use { bos ->
                                    bis.copyTo(bos)
                                }
                            }
                        }
                    } ?: appendLog("Não foi possível abrir ${childName}", "RED")
                    appendLog("Copiado: ${outFile.absolutePath}", "YELLOW")
                } catch (ex: Exception) {
                    appendLog("Erro copiando ${childName}: ${ex.message}", "RED")
                }
            }
        }
    }

    /**
     * Se o arquivo já existe, retorna um File com sufixo _1, _2, ... para evitar sobrescrita.
     */
    private fun uniqueFile(file: File): File {
        if (!file.exists()) return file
        val base = file.nameWithoutExtension
        val ext = file.extension
        var idx = 1
        while (true) {
            val newName = if (ext.isBlank()) "${base}_$idx" else "${base}_$idx.$ext"
            val candidate = File(file.parentFile, newName)
            if (!candidate.exists()) return candidate
            idx++
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        ioScope.cancel()
    }
}
