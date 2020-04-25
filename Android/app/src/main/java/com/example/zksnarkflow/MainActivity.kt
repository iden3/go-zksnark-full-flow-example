package com.example.zksnarkflow

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Button
import android.widget.TextView
import android.widget.Toast
import zkflowexample.*

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
	val zk = MobileWrapper()
        val storePath = this.filesDir.absolutePath + "/"
        val btnShow = findViewById<Button>(R.id.btnShow)
        val txt: TextView = findViewById(R.id.resTxt) as TextView
        var res = ""
        btnShow?.setOnClickListener {
            txt.text = "Running"
            try {
	    	var filesServer = "http://161.35.72.58:9000/circuit2"
		var inputs = "{\"in\":\"1\"}"
                res = zk.executeflowdownloading(storePath, filesServer, inputs)
                txt.text = res
            } catch (e: Exception) {
                Toast.makeText(this@MainActivity, "Fail: $e", Toast.LENGTH_LONG).show()
                txt.text = "Fail: $e"
            }
        }
    }
}
