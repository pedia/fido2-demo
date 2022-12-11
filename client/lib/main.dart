import 'package:flutter/material.dart';
// import 'package:url_launcher/url_launcher.dart';
import 'package:flutter_web_auth/flutter_web_auth.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'WebAuthn Demo',
      theme: ThemeData(primarySwatch: Colors.blue),
      home: const MyHomePage(title: 'WebAuthn Demo'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final controller = TextEditingController(text: 'foo');

  String? _status;

  void onRegister() {
    final name = controller.text;
    final u = Uri.parse('https://9aedu.net/register');

    final q = Map<String, List<String>>.from(u.queryParametersAll);
    q['username'] = <String>[name];
    q['next'] = <String>['foobar://success'];

    final u2 = Uri.https(u.authority, u.path, q);
    // launchUrl(u2);

    FlutterWebAuth.authenticate(url: u2.toString(), callbackUrlScheme: 'foobar')
        .then((result) {
      setState(() {
        _status = 'Got result: $result';
      });
    }).catchError((e) {
      setState(() {
        _status = 'Got error: $e';
      });
    });
  }

  void onLogin() {}

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Padding(
        padding: const EdgeInsets.all(8),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            SelectableText(
              'Ensure the default browser is Safari',
              style: Theme.of(context).textTheme.bodyLarge,
            ),
            const SizedBox(height: 24),
            SelectableText(
              _status ?? '',
              style: Theme.of(context)
                  .textTheme
                  .bodyLarge
                  ?.apply(color: Colors.red),
            ),
            const SizedBox(height: 24),
            TextField(
              controller: controller,
              decoration: const InputDecoration(labelText: 'Username'),
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 24),
            OutlinedButton(
              onPressed: onRegister,
              child: Text(
                'Register or Login',
                style: Theme.of(context).textTheme.bodyLarge,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
