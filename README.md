# DirHash
DirHash is a command line tool to take a directory and return the file hashes.

## Usage
Download DirFetch through the following command, and once DirHash is downloaded, move to the target directory using:
```
git clone https://github.com/melatonein5/DirHash.git
cd DirHash
```

#### Running through the pre-compiled binary
At the time of writing, the binary is compiled for linux systems. Run DirHash file the following `./dirhash`.
```
Usage: ./dirhash [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithm (default: md5)
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Examples:
  ./dirhash -i /path/to/dir -o output.txt -a sha256
  ./dirhash --input-dir /path/to/dir --output output.txt --algorithm sha512
  ./dirhash -t
  ./dirhash --help
```

#### Running through Go Runtime Environment (requires Go installation)
If a binary for your system is not available, you can run DirHash through the Go runtime environment with `go run dirhash.go`.
```
Usage: go run dirhash.go [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithm (default: md5)
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Examples:
  go run dirhash.go -i /path/to/dir -o output.txt -a sha256
  go run dirhash.go --input-dir /path/to/dir --output output.txt --algorithm sha512
  go run dirhash.go -t
  go run dirhash.go --help
```

#### Build from Source (requires Go installation)
To build from source, you can compile dirhash using `go build dirhash.go`, to compile an executable. The Go compiler selects your system architecture and operating system automatically. Note that if you have compiled to Windows, the `./dirhash` command becomes `dirhash.exe`.

## Example Outputs
There are two output types: terminal output; file output. Terminal output will show the results in formatted columns. File output will output to a file in csv format, regardless of file extension.

##### Terminal Output
```
[user@hostname DirHash]$ ./dirhash -a sha256
File Name                              Path                                                   Hash                                                             Hash Type
COMMIT_EDITMSG                         .git/COMMIT_EDITMSG                                    afd339433d4bff8951cd2a1a76a352adfe3bc10262399db714cbe7f0ad9eceda sha256
HEAD                                   .git/HEAD                                              1b7836f56b07eaf921c56690ef133f46ad12a3ba2cb8be9ad1cf78bb31b96995 sha256
config                                 .git/config                                            5a832d499fd863978c8ec0c3566c5fab6fe9cd7189716afa3ddb9fe7e8d3965a sha256
description                            .git/description                                       748b2e543998a28a987e49f4dc4c99b2d8f70baa04e5da2909774418edca006a sha256
applypatch-msg.sample                  .git/hooks/applypatch-msg.sample                       bcb415ff5667e3b90521aac0b09e0aae5c876ecc8c2656b4a5cd217e5e383514 sha256
commit-msg.sample                      .git/hooks/commit-msg.sample                           35bb03bd769683ee2bad41b3b0b942c66e978d8bcc396dd020e250c2f33c8840 sha256
fsmonitor-watchman.sample              .git/hooks/fsmonitor-watchman.sample                   ca9e1aac94856f82a999a65dea1b7af506a6874487e74f9226b97fab93114698 sha256
post-update.sample                     .git/hooks/post-update.sample                          062407a49e88043056ff9b23237724fdbcb4af19ccb00b5b51184b4f35ecaa8d sha256
pre-applypatch.sample                  .git/hooks/pre-applypatch.sample                       08370a3e972eba65db40247a9704dd087038c54d7794f2bf2da1bc82dc4a51a1 sha256
pre-commit.sample                      .git/hooks/pre-commit.sample                           599c01284f977cc7a5e636754ed4aceb453dcfe06cc6a5a38a1268477dffe3d6 sha256
pre-merge-commit.sample                .git/hooks/pre-merge-commit.sample                     673abbd893804293dbdc33ecfed7b481fa4ced391ce07379d356fd5d1fcb4ec5 sha256
pre-push.sample                        .git/hooks/pre-push.sample                             b4fd18fe70b172dea1da615163cf49d658172a10bd4d170ee45a554e0bd4959d sha256
pre-rebase.sample                      .git/hooks/pre-rebase.sample                           2e5ccdb3e2bad04fc5d3544e4f01ad66538008bd86363bfa4996e68a7cad4288 sha256
pre-receive.sample                     .git/hooks/pre-receive.sample                          8801d293b2b3097c7a0f7650e3f45032bda3d62823ff6a45c0b6fdcb1b75719e sha256
prepare-commit-msg.sample              .git/hooks/prepare-commit-msg.sample                   05fc76ead5fa143959bf964819653479d281b6cffa973a8e55fe6a778010e92f sha256
push-to-checkout.sample                .git/hooks/push-to-checkout.sample                     a47e51d0557ba2e578f3606c40435c454daa267c3c74fd7ed7961ed213792e54 sha256
sendemail-validate.sample              .git/hooks/sendemail-validate.sample                   c7b827489ac890b81be60d5617798ccee3f7a3985c1129d61079fd3f019f4483 sha256
update.sample                          .git/hooks/update.sample                               2a8de8c9029a70eb1dfe718c5b31c07f2ed67607cd1b5976ae1173d86582e737 sha256
index                                  .git/index                                             b070d3a25970798c78c4f7f5a5d8cf8f3bbd14f2a5b2e842fc5f577077b84662 sha256
exclude                                .git/info/exclude                                      12578d4f105793e63e8bc3853d6fdb8d51a229d0bdf923bb31690d6f3418237e sha256
HEAD                                   .git/logs/HEAD                                         ea889ef0d5f883be36f1d86593590b3c830c01c3a0eb70342e8f8acf8376bf30 sha256
main                                   .git/logs/refs/heads/main                              bdaf0a82b44bcd5c021d940eb8d18e4888bdef1e6f800fd7256c9f14699a16e2 sha256
main                                   .git/logs/refs/remotes/origin/main                     c212f430d8f6280c6c9f7607b481e5d3d7eb67426be0bd37e94d788d2ba21cdd sha256
ffd99de246994a595d09f719540338bc102824 .git/objects/2f/ffd99de246994a595d09f719540338bc102824 2592c5f4041feb3777424f7b0679f1fcd9f8074190aa5cedf5caaf9a54de1041 sha256
7f0061b697c20dc65f0abe1be1b1881d222c75 .git/objects/79/7f0061b697c20dc65f0abe1be1b1881d222c75 a6c381efe806c8e1f2a68fb2a8000fe10f824cdddf8a050234e46ade28653497 sha256
3b39e44f6d27e2f6367953794fc03abb1b89cb .git/objects/bb/3b39e44f6d27e2f6367953794fc03abb1b89cb 181e22cd1be039ab1d4ab3d7c5d469a17edcce7dce6be1de65237086441919b5 sha256
main                                   .git/refs/heads/main                                   bdea49441fc325c6bba80808e53aa6eefb4a01ddf97e99d5740e7409ccd51caf sha256
main                                   .git/refs/remotes/origin/main                          21bc691a086d8f34bda5891be3fbb67148bc6a91dbc61fa242c0b1f144fdcc9f sha256
README.md                              README.md                                              80bc0fb1a1c82259f114460ba77e427b5b3d5fcef32194ab465f3c33078a3de7 sha256
dirhash                                dirhash                                                ab08f8ca89eef473e1480ebbea862305cc10ac92ec8fe37b890d5f1b45da65bc sha256
dirhash.go                             dirhash.go                                             a45aaaf2440ee1b07e82650d2a9c7c6317ce2c7206331f643604043cc441e1e6 sha256
go.mod                                 go.mod                                                 bbbce2e1376e58c1fe3f111f2167f551229625db83fbf40abe6d03f2d8657b90 sha256
args.go                                src/args/args.go                                       66d458a380d9471c2f096ee265a8dac3cc55c06b193daa922f3cbccc63b8d570 sha256
hash_algorithm.go                      src/args/hash_algorithm.go                             f7b5755e3a7f6f9e4de2f4405853420e4aabffba4158df7623907a183d477d82 sha256
parse_args.go                          src/args/parse_args.go                                 d4c99f7cba39c251430b3983c85faf2b26d908b78e53faa12b28b458b115fb04 sha256
help.go                                src/cmdline/help.go                                    95cc5040d7f297c061c0482e5e25e2e241f9a4e2a036b14229e40772c3ec22f0 sha256
output_files.go                        src/cmdline/output_files.go                            a19e0735bcf3175290476c7a48be3a0f55345341d5fe3d76dace46eb0a488a3c sha256
enumerate_files.go                     src/files/enumerate_files.go                           5660216c3957bf5f2bfeaed4827beaaef35ab80392fc38159931bb31675da312 sha256
file.go                                src/files/file.go                                      791f770d5ee7b622b668dc720d22dcb1af9a793b57f0a17278c314d59f7c62db sha256
hash_files.go                          src/files/hash_files.go                                5847d97a7f2c68f27555aad5c05bfd12a7bfe3963d86785edd4223434284a7b6 sha256
md5_files.go                           src/files/md5_files.go                                 18a1bef74d7ac79e32a0ef111e7d615f2311fbefe2d3f96eb6852a6f3a7446be sha256
sha1_files.go                          src/files/sha1_files.go                                20fa18cb29b5009cd992114aea2866a80d3b5fdd4de5826497a63270dd5dadff sha256
sha256_files.go                        src/files/sha256_files.go                              9786814c4d3bb437956591d0a20914a69f4cfbde8233bb4d5ab51078685cbeee sha256
sha512_files.go                        src/files/sha512_files.go                              4dfea9664efd8f59d22421d21d712c4a916cf294a638d32b06f7d8d9cfbaa3f0 sha256
write_output.go                        src/files/write_output.go                              0eac73e2af1c90940d91b70c369fc9a183e4dba263e06b0f498db91f9d2b40f2 sha256
```

#### CSV Output
```
Path,Hash,HashType
.git/COMMIT_EDITMSG,34bb1ef88028a1b61d48100cc9e1dfa37b8af19f1e3babe2747d4f03dbb9d4ac4f464cb01e8e11817e6b2050661a6d3f6c8da21bd5608e236bfcd3f87083b9a7,sha512
.git/HEAD,d2af064a992de738fe2445adddeb0fb6085e74a698937dacaed294ae29376c869bdaa9c0ca1453f8014c168368916dcdf851a595b8515515d7f8f750fbd2acdb,sha512
.git/config,686dee4a9d02132ad1d5f50843952685a330d8e9f6b65fd21cd3828d690695f61f88597cb4b913ec8d8673e5349dd7fefa6d91c2ae8939c28449810bab7e5747,sha512
.git/description,6b7a652443e815dee8d16e71ac80625fe78affeaf54ba190ac229950ba9072bd89c90a5d2a01224070c827358ff185f57a2efe7abd2b1333cb5da84102c2e784,sha512
.git/hooks/applypatch-msg.sample,9f2d3ac648424b9c6a2e462837ef0666d4d89820921fec173a55a003c941b63b2601cae0eb57938b2056ccb571d5df8bbfd2617bb3a8bd1c8f0c8607beb2b837,sha512
.git/hooks/commit-msg.sample,8d5395d45e8718a94ffa766b9cdbf064ef31557fdbd22cffb090fb9ccf41268b837b5efa84742272bf2c95e97fd1a7146a6a3cf00de86caf09440fa63bae84ec,sha512
.git/hooks/fsmonitor-watchman.sample,c1ace3bc31e0fe3076f0b97d7fe89c416f339b894f925539a7c358d697ddc3e3d2e25740c8c2da41018f16e419f954e67adbd4bc2129ecef5b798b951a6d0a62,sha512
.git/hooks/post-update.sample,73ca6bfcb2e478fe87cc2f1e7dff55bb3f2588d7785b3b4bc180c83b75ad4cb456c6d3fdd370b52a1a4887f250db51ea7ceac5fad12a8e519f3b5ebcfc77a6f4,sha512
.git/hooks/pre-applypatch.sample,22e9e24f8343baef07bcfe7e412b3c67755f3d3eb5445193419340091f72575722e4a906a5aa2115f44aa169e64fbf5c237f6ab5cf6fba054af4a7252d7ae5d2,sha512
.git/hooks/pre-commit.sample,8587e03cd0a452bc917efc34af72c2876be7d4adf7347cdc4db37338ee283b0d3b04b4a111be430bc851e0b27888e2e96a00807a368d77730e9db5c1a4af166c,sha512
.git/hooks/pre-merge-commit.sample,8422a1c556ba91a12848471aeb928ed51c3d77f03e3e3b71cd637b40dca25ca4f23c4ffe51e8d099dba9b09c368aff0577fe7e23feb3e1fccab1cef27029a39b,sha512
.git/hooks/pre-push.sample,8ad5fdf9fd482f5909e8375fc647f0a91e1b15de83e427609127632b9b9b892c3d5c4bacf8280ebba1c94d589fc38fead3d74145862514a8326223da70ef51b7,sha512
.git/hooks/pre-rebase.sample,11ed8d5a3ebdfea1ec098aae2a39dda4bf673f3e04a841a798ed4dbc76f628ae293b2bab8ec3a493aa28e4d822fa8e9d56fafc9aad8b285d26d2923b33ce827d,sha512
.git/hooks/pre-receive.sample,6ff196f0f66524e96cb11200fdee957051de399d0519c3ec1d1d4826bdd09fe4a9bf81ac6201c6c266bf30eea2c1d069ae2cb4c504e0dad9aee0055852243a8b,sha512
.git/hooks/prepare-commit-msg.sample,a3e4495f2e707e714f3e3342b3d90a695b1f579ae2f442c057f415f18b25935a67a46e9503ea007476e23eaf3b1a4a7a8d7b725aadd4d0b8d613ff413521ea0d,sha512
.git/hooks/push-to-checkout.sample,207eb72a68379db76dc6a01440448bfc4b3c5ac5f1ab31ab5179bad9bc6adb066a275e892dc6d5cca020564345f512efdea11d008120ef5853185ca1beecf526,sha512
.git/hooks/sendemail-validate.sample,a858731c5fc2e87a7155d4f9f4be0730af763e7d8e71fac154f5303e2ded716308728dd2583449461a40688eeff618cf62dd388ff61ac61ce9d7daebc9ca8698,sha512
.git/hooks/update.sample,a1910b35e0af53a3b5bba0872ae2198d55962f2611692f0dc0ada15f6789f674655673810c065c869f86be723c2b0a33a6efc7d5ebeac2cd1abcdbc99badc737,sha512
.git/index,f3a27d665ec7d39db84e7b7bf9621269240a0609c6c7d7987dfdc870f83b2c6d786ec1781bf0bb4d73107c94ba1144d311bb875bdb09629da1d609d132b48675,sha512
.git/info/exclude,0f018163e053ef1c6e75f00fee2ac08e38795104722574d85901a3e05d73d1df17714a6053e834b8936d2dd69c7b435e171b5ccffbfa04d39c04810ac1522208,sha512
.git/logs/HEAD,f6d2b6cbf9e4f5c8440c2b4934b25f997807d181bd09fc18d9db3b9a65e827ff2effe09bc3940132c1f470d5129054ecb8833bb21b800c9e80fe8468f6ba380b,sha512
.git/logs/refs/heads/main,11fbbc6c622a1618edf5a8814fbf7707fab2174d7b3c34f2d311df06877ebea684858f9204e52889b3282359ee6f955f4859295d9cdadb925723debd045f8489,sha512
.git/logs/refs/remotes/origin/main,ed8a3510ee878b1e37793b9a1386e61e97fcf583895bcb1f96b74ac662a67a3619cb3a438a38671124278d8556f21e473adf4331afe9167568105ccfad7cca77,sha512
.git/objects/2f/ffd99de246994a595d09f719540338bc102824,fca87af0078af2d48e7c42d423af6c191ad7180a337a99d8b493cd2ba3dfcb3bb8cf2b7be8329580d53a3392b412a71a99905b953183a348387a8fbe3cffdd7f,sha512
.git/objects/79/7f0061b697c20dc65f0abe1be1b1881d222c75,9892bd79feb79d1e5f051aaba40a7c47936e88ce400bce14d7708ce2016bff8b9decd00988be62a7509cfc3d9507df09fbbc099b3c7664bed8982558bd6be795,sha512
.git/objects/bb/3b39e44f6d27e2f6367953794fc03abb1b89cb,1a9a7b4ab724b3a2937213136a22e4f03d9150e72231749df1c1a53cf62fdcd130727c066c9173f6e81efd5e3b2de0eb7defdc688affd3f12100a98753bd23aa,sha512
.git/refs/heads/main,fc935d2217d4f30f4d55bde62503063330b5846164c7c8f1934355e9173b2a585cce2d094e9442577fb0600dea53dc2e458d547c70889e088238c974c4823134,sha512
.git/refs/remotes/origin/main,82c00459623eb0939e36ce8b34f7a15b4cf14da47bc561d65f6757178fad9fd5a300da29209da99cf0993d3bc9446ffa69ab7fefd92ac77eb1afba0e955141bc,sha512
README.md,827d8ea799490cb1d8721239a9dd719e0d6de3afb661a724f601e8ee4f99b399a2216da83e432a5b16aae6abb3b88aa8c32a5226a49b457b098c248f3e157588,sha512
dirhash,9f63843a3082cbaea927b162e40241a7a91990fd25ddb6ffb6667964abe82a09b7d7bfdcb0ff813b712cf5b29fb93c305924badf9b06d52a9093673290f1b1d2,sha512
dirhash.go,6f5339d4f462a0a363d58c05afe747937deb5cffbb9839e6a2cc9431178f104e45ca27118101a34fe086e0a7001fb58d155251bf52d6bd4c27472544f2b94cf3,sha512
go.mod,ad6773ab853200b272eec31380185a41d9f832878fd241764756b9a68cb71f971d1a80271fb9ac3a9f596e3ae9620c87c995fe84423eec8e64afebe2a0ad2f7d,sha512
src/args/args.go,42772e9f5c9730dad50ef04bb086d2f29d8974bea9946d51bd1401fe39a15151ae5a1a98778542d35b4ff9fd83f856fb86ddaa5ac19f7d683a0df015e5b0f793,sha512
src/args/hash_algorithm.go,519049bd042054c3a1158d27fc259fa1d53feac68c8a1db837c8f9049dd4257d5d08cdba997bfc05ec804aa01e221538d04db4c7c4e431aba19c4eb5666af2ab,sha512
src/args/parse_args.go,5245d2ec0c3ec77fc97d9a89b1a1a6b10df64aea82646790b5cff925bbb9ca8bde8ed3df7c9c505727ab0f0c36479c30f4db081fbd0acaccd40f691681c596ac,sha512
src/cmdline/help.go,9a0d3d9bd9e73dc4855b78f8ff8dd7db1c623a07563b404a976071097a6e4ba39c39d1e08c526fec1f009b836808951555c37ff90f83237b14bff56601427acc,sha512
src/cmdline/output_files.go,3dd2c94e2b90586973d8d65169efbceeb730b28faae7510bfd63d0f8a24cab76cefc27279d9a428e9c797939f40f723f63f416ee51002c03105ce4afe341efe2,sha512
src/files/enumerate_files.go,232aa8c8a49ce73db529baede71bb9e3c401693a481171902f6406d27f3fcfe1a260ac9316802830f94cb308c4373237bd56d718c9fa3c9f8ce59d69b5036529,sha512
src/files/file.go,670ba48849464991f8070ebfcffd2bb5a86781faae38492cab56331be9803a549a9064ee4def4865872bb545eb7a2ace4d2d94e8f3495c9f9db765d62a6f3d4e,sha512
src/files/hash_files.go,6865f9a842998a822dd4ed481c86a2d0081e37998d0bf1681ed0133f333d916330cd7a47a6eaf1c78dbddcdd083164125191be240c620da2e260ac0a82b19416,sha512
src/files/md5_files.go,e03534dd45d373e743a65bf622bd3815286f101092c713f17df53baf4341259165969ee4f9409e5dd2fe64cd55660dd926fe069bb34bfa4442fff77a5555b729,sha512
src/files/sha1_files.go,e66eb9f2a87e21051f46b6590f6e7cd5572fa35b74024a35728a36f8c61d2421b9155f3cbec67fcc3e87bed26029f6418b218aeefa022446af35b1727eb3f32f,sha512
src/files/sha256_files.go,182f1a84ad40d3c4dc50aa9dc81e7e316aae6172a7a66621052765d296ba0ae889fd4c3c59e60d784754f2cdd39d6bd320fa52e420c98ce0a2816fc87cd36b0a,sha512
src/files/sha512_files.go,b03f1a3d9e3f64f11f81b4011ff3579f7b8a93ea97c87ac1b5544fe9e5e03b054cbc37650aadec7b93e31f2e29ceb750b89fd268c5a424188bb1a27235bcab08,sha512
src/files/write_output.go,3eb5bb989a47bfc677027578be94f8d3190014c182839aa5eb7cb08a6ad9ddb3c29926d95c9616e1a58ebc9accacf9a3463c36ea844dd22342437b053c207008,sha512
```

## Roadmap
The three main improvements to make are: 
1. Running multiple hash functions in a single command.
1. Add installation scripts to allow system-wide access at the terminal.
1. Support for more hashing functions.