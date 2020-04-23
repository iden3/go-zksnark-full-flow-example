package com.example.zksnarkflow

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Button
import android.widget.Toast
import zkflowexample.*

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        val zk = MobileZKFlow()
        val storePath = this.filesDir.absolutePath + "/"
        val btnShow = findViewById<Button>(R.id.btnShow)
        btnShow?.setOnClickListener {
            try {
                zk.run(storePath)
                Toast.makeText(this@MainActivity, "Flow completed!", Toast.LENGTH_LONG).show()
            } catch (e: Exception) {
                Toast.makeText(this@MainActivity, "Fail: $e", Toast.LENGTH_LONG).show()
            }
        }
    }
}
